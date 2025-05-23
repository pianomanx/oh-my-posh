package segments

import (
	"errors"
	"testing"
	"time"

	"github.com/jandedobbeleer/oh-my-posh/src/properties"
	"github.com/jandedobbeleer/oh-my-posh/src/runtime/mock"

	"github.com/stretchr/testify/assert"
	testify_ "github.com/stretchr/testify/mock"
)

type mockedWithingsAPI struct {
	testify_.Mock
}

func (s *mockedWithingsAPI) GetMeasures(meastypes string) (*WithingsData, error) {
	args := s.Called(meastypes)
	return args.Get(0).(*WithingsData), args.Error(1)
}

func (s *mockedWithingsAPI) GetActivities(activities string) (*WithingsData, error) {
	args := s.Called(activities)
	return args.Get(0).(*WithingsData), args.Error(1)
}

func (s *mockedWithingsAPI) GetSleep() (*WithingsData, error) {
	args := s.Called()
	return args.Get(0).(*WithingsData), args.Error(1)
}

func TestWithingsSegment(t *testing.T) {
	cases := []struct {
		MeasuresError   error
		ActivitiesError error
		SleepError      error
		WithingsData    *WithingsData
		Case            string
		ExpectedString  string
		Template        string
		ExpectedEnabled bool
	}{
		{
			Case:            "Error",
			MeasuresError:   errors.New("error"),
			ActivitiesError: errors.New("error"),
			SleepError:      errors.New("error"),
			ExpectedEnabled: false,
		},
		{
			Case: "Only Measures data",
			WithingsData: &WithingsData{
				Body: &Body{
					MeasureGroups: []*MeasureGroup{
						{
							Measures: []*Measure{
								{
									Value: 7077,
									Unit:  -2,
								},
							},
						},
					},
				},
			},
			ActivitiesError: errors.New("error"),
			SleepError:      errors.New("error"),
			ExpectedEnabled: true,
			ExpectedString:  "70.77kg",
		},
		{
			Case: "Multiple Measuring Groups, only Measures data",
			WithingsData: &WithingsData{
				Body: &Body{
					MeasureGroups: []*MeasureGroup{
						{
							Measures: []*Measure{
								{
									Value: 7123,
									Unit:  -2,
								},
							},
						},
						{
							Measures: []*Measure{
								{
									Value: 7754,
									Unit:  -2,
								},
							},
						},
					},
				},
			},
			ActivitiesError: errors.New("error"),
			SleepError:      errors.New("error"),
			ExpectedEnabled: true,
			ExpectedString:  "77.54kg",
		},
		{
			Case: "Measures, no data",
			WithingsData: &WithingsData{
				Body: &Body{},
			},
			ActivitiesError: errors.New("error"),
			SleepError:      errors.New("error"),
			ExpectedEnabled: false,
		},
		{
			Case:           "Activities",
			Template:       "{{ .Steps }} steps",
			ExpectedString: "7077 steps",
			WithingsData: &WithingsData{
				Body: &Body{
					Activities: []*Activity{
						{
							Steps: 5066,
							Date:  time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
						},
						{
							Steps: 7077,
							Date:  time.Now().Format("2006-01-02"),
						},
					},
				},
			},
			MeasuresError:   errors.New("error"),
			SleepError:      errors.New("error"),
			ExpectedEnabled: true,
		},
		{
			Case:           "Sleep",
			Template:       "{{ .SleepHours }}hr",
			ExpectedString: "11.8hr",
			WithingsData: &WithingsData{
				Body: &Body{
					Series: []*Series{
						{
							Startdate: 1594159200,
							Enddate:   1594201500,
						},
					},
				},
			},
			MeasuresError:   errors.New("error"),
			ActivitiesError: errors.New("error"),
			ExpectedEnabled: true,
		},
		{
			Case:           "Sleep and Activity",
			Template:       "{{ .Steps }} steps with {{ .SleepHours }}hr of sleep",
			ExpectedString: "976 steps with 11.8hr of sleep",
			WithingsData: &WithingsData{
				Body: &Body{
					Series: []*Series{
						{
							Startdate: 1594159200,
							Enddate:   1594201500,
						},
					},
					Activities: []*Activity{
						{
							Steps: 976,
							Date:  time.Now().Format("2006-01-02"),
						},
					},
				},
			},
			MeasuresError:   errors.New("error"),
			ExpectedEnabled: true,
		},
	}

	for _, tc := range cases {
		api := &mockedWithingsAPI{}
		api.On("GetMeasures", "1").Return(tc.WithingsData, tc.MeasuresError)
		api.On("GetActivities", "steps").Return(tc.WithingsData, tc.ActivitiesError)
		api.On("GetSleep").Return(tc.WithingsData, tc.SleepError)

		withings := &Withings{
			api: api,
		}
		withings.Init(properties.Map{}, &mock.Environment{})

		enabled := withings.Enabled()
		assert.Equal(t, tc.ExpectedEnabled, enabled, tc.Case)
		if !enabled {
			continue
		}

		if tc.Template == "" {
			tc.Template = withings.Template()
		}

		var got = renderTemplate(&mock.Environment{}, tc.Template, withings)
		assert.Equal(t, tc.ExpectedString, got, tc.Case)
	}
}
