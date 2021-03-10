package profitbricks

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
)

type VolumeSuite struct {
	ClientBaseSuite
}

func TestVolume(t *testing.T) {
	suite.Run(t, new(VolumeSuite))
}

func (s *VolumeSuite) Test_CreateSnapshot() {
	mRsp := makeJsonResponse(http.StatusAccepted, loadTestData(s.T(), "create_snapshot.json"))
	mRsp.Header.Set("Location", "requests/foo/status")
	httpmock.RegisterResponder(http.MethodPost,
		"=~/datacenters/111/volumes/222/create-snapshot", httpmock.ResponderFromResponse(mRsp))
	rsp, err := s.c.CreateSnapshot("111", "222", "TestSnapshot01", "")
	s.NoError(err)
	s.Equal("requests/foo/status", rsp.Headers.Get("location"))
	s.Equal(float32(0.0), rsp.Properties.Size)
}

func (s *VolumeSuite) Test_CreateSnapshotAndWait() {
	mRsp := makeJsonResponse(http.StatusAccepted, loadTestData(s.T(), "create_snapshot.json"))
	mRsp.Header.Set("Location", "requests/foo/status")
	httpmock.RegisterResponder(http.MethodPost,
		"=~/datacenters/111/volumes/222/create-snapshot", httpmock.ResponderFromResponse(mRsp))
	httpmock.RegisterResponder(http.MethodGet,
		"=~/requests/foo/status", httpmock.ResponderFromResponse(
			makeJsonResponse(http.StatusOK, loadTestData(s.T(), "request_create_snapshot_done.json")),
		))
	httpmock.RegisterResponder(http.MethodGet,
		"=~/snapshots/9676a76c-e2a0-4365-af4d-8ab3fa8bf997", httpmock.ResponderFromResponse(
			makeJsonResponse(http.StatusOK, loadTestData(s.T(), "get_snapshot.json")),
		))

	rsp, err := s.c.CreateSnapshotAndWait(context.Background(), "111", "222", "TestSnapshot01", "")
	s.NoError(err)
	s.Equal(float32(100.0), rsp.Properties.Size)
}
