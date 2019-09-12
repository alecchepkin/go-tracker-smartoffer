package faker

import (
	"errors"
	"testing"
)

type mockStatRepo struct {
	stats []stat
	clk   struct {
		click
		error
	}
	receivedFindSmartStats bool
	receivedFindAnyClick   bool
}

func (s *mockStatRepo) findSmartStats() []stat {
	s.receivedFindSmartStats = true
	return s.stats
}

func (s *mockStatRepo) findAnyClick(offerId int, publisherId int) (clk click, err error) {
	s.receivedFindAnyClick = true
	return s.clk.click, s.clk.error
}

func newMockStatRepo(stats []stat, clk click, err error) mockStatRepo {
	return mockStatRepo{
		stats: stats,
		clk: struct {
			click
			error
		}{click: clk, error: err},
	}
}

type mockFakeRepo struct {
	receivedSave bool
}

func (s *mockFakeRepo) save(f fake) {
	s.receivedSave = true
}

func TestFakeCreator_Process(t *testing.T) {

	type fields struct {
		statRepo mockStatRepo
	}

	type wants struct {
		receivedFindSmartStats bool
		receivedFindAnyClick   bool
		receivedSave           bool
	}
	tests := []struct {
		name   string
		fields fields
		wants  wants
	}{
		{"should pass at the end", fields{statRepo: newMockStatRepo([]stat{{1, 1, clicksInterval + 1, 0}}, click{}, nil)},
			wants{true, true, true}},
		{"should not save if click not found", fields{statRepo: newMockStatRepo([]stat{{1, 1, clicksInterval + 1, 0}}, click{}, errors.New("foo"))},
			wants{true, true, false}},
		{"should not finding click if should not create fake", fields{statRepo: newMockStatRepo([]stat{{1, 1, 0, maxFakes}}, click{}, errors.New("foo"))},
			wants{true, false, false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr := mockFakeRepo{}
			fc := FakeCreator{
				statRepo: &tt.fields.statRepo,
				fakeRepo: &fr,
			}
			fc.Process()
			if tt.fields.statRepo.receivedFindSmartStats != tt.wants.receivedFindSmartStats {
				t.Errorf("receivedFindSmartStats = %t, want %t", tt.fields.statRepo.receivedFindSmartStats, tt.wants.receivedFindSmartStats)
			}
			if tt.fields.statRepo.receivedFindAnyClick != tt.wants.receivedFindAnyClick {
				t.Errorf("receivedFindAnyClick = %t, want %t", tt.fields.statRepo.receivedFindAnyClick, tt.wants.receivedFindAnyClick)
			}
			if fr.receivedSave != tt.wants.receivedSave {
				t.Errorf("receivedSave = %t, want %t", fr.receivedSave, tt.wants.receivedSave)
			}
		})
	}
}

func Test_needFake(t *testing.T) {
	type args struct {
		stat stat
	}
	tests := []struct {
		name    string
		args    args
		wantRsl bool
	}{
		{"enough fakes if max", args{stat: stat{clicks: clicksInterval, convs: maxFakes}}, false},
		{"enough fakes if max (more clicks)", args{stat: stat{clicks: clicksInterval * 10, convs: maxFakes}}, false},
		{"need first fake", args{stat: stat{clicks: clicksInterval + 1, convs: 0}}, true},
		{"need next fake ", args{stat: stat{clicks: clicksInterval*2 + 1, convs: maxFakes - 1}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRsl := needFake(tt.args.stat); gotRsl != tt.wantRsl {
				t.Errorf("needFake() = %v, want %v", gotRsl, tt.wantRsl)
			}
		})
	}
}
