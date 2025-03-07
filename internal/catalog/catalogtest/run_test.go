package catalogtest

import (
	"testing"

	svctest "github.com/shulutkov/yellow-pages/agent/grpc-external/services/resource/testing"
	"github.com/shulutkov/yellow-pages/internal/catalog"
	"github.com/shulutkov/yellow-pages/internal/catalog/internal/controllers"
	"github.com/shulutkov/yellow-pages/internal/controller"
	"github.com/shulutkov/yellow-pages/internal/resource/reaper"
	"github.com/shulutkov/yellow-pages/proto-public/pbresource"
	"github.com/shulutkov/yellow-pages/sdk/testutil"
)

func runInMemResourceServiceAndControllers(t *testing.T, deps controllers.Dependencies) pbresource.ResourceServiceClient {
	t.Helper()

	ctx := testutil.TestContext(t)

	// Create the in-mem resource service
	client := svctest.RunResourceService(t, catalog.RegisterTypes)

	// Setup/Run the controller manager
	mgr := controller.NewManager(client, testutil.Logger(t))
	catalog.RegisterControllers(mgr, deps)

	// We also depend on the reaper to take care of cleaning up owned health statuses and
	// service endpoints so we must enable that controller as well
	reaper.RegisterControllers(mgr)
	mgr.SetRaftLeader(true)
	go mgr.Run(ctx)

	return client
}

func TestControllers_Integration(t *testing.T) {
	client := runInMemResourceServiceAndControllers(t, catalog.DefaultControllerDependencies())
	RunCatalogV1Alpha1IntegrationTest(t, client)
}
