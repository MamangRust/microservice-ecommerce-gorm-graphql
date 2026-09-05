package banner_test

import (
	"net/http"
	"strconv"
	"testing"

	graphtest "github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type BannerApiTestSuite struct {
	tests.BaseTestSuite
	handler  http.Handler
	bannerID int
}

func (s *BannerApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupBannerService()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupAuthService()
	s.SetupMerchantService()
	s.SetupOrderService()
	s.SetupTransactionService()

	resolver := graphtest.NewResolver(s.Conns, s.Log, s.GetCacheStore())
	s.handler = graphtest.NewHandler(resolver)

	// Seed a banner via DB
	err := s.DBPool().QueryRow(s.Ctx,
		`INSERT INTO banners (name, start_date, end_date, start_time, end_time, is_active) VALUES ('Seed Banner', '2026-01-01', '2026-12-31', '08:00:00', '20:00:00', true) RETURNING banner_id`,
	).Scan(&s.bannerID)
	s.Require().NoError(err)
}

func (s *BannerApiTestSuite) TestBannerApiLifecycle() {
	idStr := strconv.Itoa(s.bannerID)

	// 1. FindById
	resp, err := graphtest.ExecuteGraphQL(s.handler, `query { findBannerById(input: { id: `+idStr+` }) { status message data { banner_id name } } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findBannerById"].(map[string]interface{})["status"])

	// 2. FindAll
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findAllBanners(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findAllBanners"].(map[string]interface{})["status"])

	// 3. FindByActive
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findActiveBanners(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findActiveBanners"].(map[string]interface{})["status"])

	// 4. Update
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { updateBanner(input: { banner_id: `+idStr+`, name: "Updated Banner", start_date: "2026-01-01", end_date: "2026-12-31", start_time: "08:00:00", end_time: "20:00:00", is_active: true }) { status message data { banner_id name } } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["updateBanner"].(map[string]interface{})["status"])

	// 5. Trash
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { trashBanner(input: { id: `+idStr+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["trashBanner"].(map[string]interface{})["status"])

	// 6. FindByTrashed
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findTrashedBanners(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findTrashedBanners"].(map[string]interface{})["status"])

	// 7. Restore
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { restoreBanner(input: { id: `+idStr+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["restoreBanner"].(map[string]interface{})["status"])

	// 8. Re-trash + DeletePermanent
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { trashBanner(input: { id: `+idStr+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["trashBanner"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { deleteBannerPermanent(input: { id: `+idStr+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["deleteBannerPermanent"].(map[string]interface{})["status"])

	// 9. RestoreAll
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { restoreAllBanners { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["restoreAllBanners"].(map[string]interface{})["status"])

	// 10. DeleteAll
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { deleteAllBannersPermanent { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["deleteAllBannersPermanent"].(map[string]interface{})["status"])
}

func TestBannerApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(BannerApiTestSuite))
}
