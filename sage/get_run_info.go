package sage

import (
	"context"
	"encoding/json"
	"fmt"

	"strings"

	"github.com/ws6/fcinfo"

	"github.com/araddon/dateparse"
	"github.com/ws6/bsshgo"
)

//GetMsi a GET method with a return type map[string]interface{}
func GetMsi(ctx context.Context, client *bsshgo.Client, url string) (map[string]interface{}, error) {

	body, err := client.GetBytes(ctx, url)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]interface{})
	if err := json.Unmarshal(body, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}

type SeqStatsPre struct {
	Run struct {
		Id           string
		Name         string
		DateCreated  string
		DateModified string
		Href         string
	}
	Reads []struct {
		ReadNumber               int
		IsIndexed                bool
		TotalCycles              int
		YieldTotal               float64
		ProjectedTotalYield      float64
		PercentAligned           float64
		ErrorRate                float64
		IntensityCycle1          float64
		PercentGtQ30             float64
		PercentGtQ30Last10Cycles float64
	}
	Lanes []struct {
		Id                       string //lane Id
		Href                     string
		Density                  int64
		ErrorRate                float64
		ErrorRate100             float64
		ErrorRate35              float64
		ErrorRate50              float64
		ErrorRate75              float64
		IntensityCycle1          float64
		LaneNumber               int
		MaxCycleCalled           int
		MaxProjectedYieldInGbp   float64
		PercentAligned           float64
		PercentGtQ30             float64
		PercentGtQ30Last10Cycles float64
		PercentPf                float64
		Phasing                  float64
		PrePhasing               float64
		ProjectedYieldInGbp      float64
		Reads                    int64
		ReadsPf                  int64
		Status                   string
		Yield                    float64
		PhasingSlope             float64
		PhasingOffset            float64
		PrePhasingSlope          float64
		PrePhasingOffset         float64
	}
	LanesByRead []struct {
		ReadNumber               int
		LaneNumber               int
		TileCount                int
		Density                  int64
		DensityDeviation         float64
		PercentPf                float64
		PercentPfDeviation       float64
		Phasing                  float64
		PrePhasing               float64
		Reads                    int64
		ReadsPf                  int64
		PercentGtQ30             float64
		PercentGtQ30Last10Cycles float64
		Yield                    float64
		MinCycleCalled           int
		MaxCycleCalled           int
		MinCycleError            float64
		MaxCycleError            float64
		PercentAligned           float64
		PercentAlignedDeviation  float64
		ErrorRate                float64

		ErrorRateDeviation       float64
		ErrorRate35              float64
		ErrorRate35Deviation     float64
		ErrorRate50              float64
		ErrorRate50Deviation     float64
		ErrorRate75              float64
		ErrorRate75Deviation     float64
		ErrorRate100             float64
		ErrorRate100Deviation    float64
		IntensityCycle1          float64
		IntensityCycle1Deviation float64
		PhasingSlope             float64
		PhasingOffset            float64
		PrePhasingSlope          float64
		PrePhasingOffset         float64
		PercentNoCalls           float64
		ClusterDensity           int64
		Occupancy                float64
	}
	Chemistry                     string
	ErrorRate                     float64
	ErrorRateR1                   float64
	ErrorRateR2                   float64
	Href                          string
	IntensityCycle1               float64
	IsIndexed                     bool
	MaxCycleCalled                int
	MaxCycleExtracted             int
	MaxCycleScored                int
	MinCycleCalled                int
	MinCycleExtracted             int
	MinCycleScored                int
	NonIndexedErrorRate           float64
	NonIndexedIntensityCycle1     float64
	NonIndexedPercentAligned      float64
	NonIndexedPercentGtQ30        float64
	NonIndexedProjectedTotalYield float64
	NonIndexedYieldTotal          float64
	NumCyclesIndex1               int
	NumCyclesIndex2               int
	NumCyclesRead1                int
	NumCyclesRead2                int
	NumLanes                      int
	NumReads                      int
	NumSurfaces                   int
	NumSwathsPerLane              int
	NumTilesPerSwath              int
	PercentAligned                float64
	PercentGtQ30                  float64
	PercentGtQ30R1                float64
	PercentGtQ30R2                float64
	PercentGtQ30Last10Cycles      float64
	PercentPf                     float64
	PercentResynthesis            float64
	PhasingR1                     float64
	PhasingR2                     float64
	PrePhasingR1                  float64
	PrePhasingR2                  float64
	ProjectedTotalYield           float64
	ReadsPfTotal                  int64
	ReadsTotal                    int64
	YieldTotal                    float64
	Clusters                      int64
	ClustersPf                    int64
	ClusterDensity                int64
	Occupancy                     float64
}

func GetSequencingStats(ctx context.Context, client *bsshgo.Client, runId string) (*SeqStatsPre, error) {
	url := fmt.Sprintf(`/v2/runs/%s/sequencingstats`, runId)
	body, err := client.GetBytes(ctx, url)
	if err != nil {
		return nil, err
	}
	ret := new(SeqStatsPre)
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

//basic runInfo data type
type RunInfoPre struct {
	Id             string
	Name           string
	ExperimentName string
	DateCreated    string
	DateModified   string
	Status         string
	Instrument     struct {
		Id           int64
		Name         string
		Number       int64
		Type         string
		PlatformName string
	}
	InstrumentRunStatus     string
	DateInstrumentStarted   string
	DateInstrumentCompleted string
	FlowcellBarcode         string
	ReagentBarcode          string
	FlowcellPosition        string
	LaneAndQcStatus         string
	Workflow                string
	SampleSheetName         string
	TotalSize               int64
	UploadStatus            string
	DateUploadStarted       string
	DateUploadCompleted     string
	IsArchived              bool
	Href                    string //can be used as flowcell.location equivelent
	V1Pre3Id                string
}

func GetRunInfo(ctx context.Context, client *bsshgo.Client, runId string) (*RunInfoPre, error) {
	url := fmt.Sprintf(`/v2/runs/%s`, runId)
	body, err := client.GetBytes(ctx, url)
	if err != nil {
		return nil, err
	}
	ret := new(RunInfoPre)
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func toSageFlowcellStatus(s string) string {
	if s == `Complete` {
		return `flowcell.finished`
	}
	return `flowcell.started`
}

func getReadLength(seq *SeqStatsPre) string {

	ret := []string{}
	if seq.NumCyclesRead1 != 0 {
		ret = append(ret, fmt.Sprintf("%d", seq.NumCyclesRead1))
	}
	if seq.NumCyclesIndex1 != 0 {
		ret = append(ret, fmt.Sprintf("%d", seq.NumCyclesIndex1))
	}
	if seq.NumCyclesIndex2 != 0 {
		ret = append(ret, fmt.Sprintf("%d", seq.NumCyclesIndex2))
	}
	if seq.NumCyclesRead2 != 0 {
		ret = append(ret, fmt.Sprintf("%d", seq.NumCyclesRead2))
	}

	return strings.Join(ret, ",")
}
func getCycles(seq *SeqStatsPre) int {
	return seq.NumCyclesRead1 + seq.NumCyclesIndex1 + seq.NumCyclesIndex2 + seq.NumCyclesRead2

}

func BuildSageSpec(seq *SeqStatsPre, run *RunInfoPre, fc *fcinfo.Flowcell) (map[string]interface{}, error) {
	ret, err := _buildSageSpec(seq, run)
	if err != nil {
		return nil, err
	}

	if fc == nil {
		return ret, nil
	}
	flowcell, ok := ret[`flowcell`].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf(`not msi`)
	}

	flowcell[`flowcell_barcode`] = fc.FlowcellBarcode
	if fc.RunParamOutputFolder != "" {
		flowcell[`run_param_output_folder`] = fc.RunParamOutputFolder
	}
	if fc.RunId != "" {
		flowcell[`run_id`] = fc.RunId
	}

	flowcell[`chemistry`] = fc.Chemistry
	flowcell[`machine_name`] = fc.MachineName
	if t, err := dateparse.ParseLocal(fc.RunStartDate); err == nil {
		if !t.IsZero() {
			flowcell[`run_start_date`] = t.String()
		}
	}

	if run.Instrument.Number != 0 {

		flowcell[`run_number`] = run.Instrument.Number

	}

	if run.Instrument.Type != "" {
		flowcell[`instrument_type`] = run.Instrument.Type
	}

	if fc.ApplicationVersion != "" {
		flowcell[`application_version`] = fc.ApplicationVersion
	}
	if fc.ApplicationName != "" {
		flowcell[`application_name`] = fc.ApplicationName
	}

	if fc.FpgaVersion != "" {
		flowcell[`fpga_version`] = fc.FpgaVersion
	}

	if fc.RtaVersion != "" {
		flowcell[`rta_version`] = fc.RtaVersion
	}

	return ret, nil
}

// BuildSageSpec convert *SeqStatsPre, *RunInfoPre to below  (sage use)specs
//https://confluence.illumina.com/display/FBS/sage+incoming+specs
func _buildSageSpec(seq *SeqStatsPre, run *RunInfoPre) (map[string]interface{}, error) {
	ret := make(map[string]interface{})
	//add flowcell

	flowcell := make(map[string]interface{})
	ret[`flowcell`] = flowcell
	flowcell[`location`] = run.Href

	flowcell[`run_id`] = run.Name
	flowcell[`status`] = toSageFlowcellStatus(run.Status)
	flowcell[`indexed`] = `no`
	if seq.IsIndexed {
		flowcell[`indexed`] = `yes`
	}
	flowcell[`flowcell_barcode`] = run.FlowcellBarcode
	flowcell[`read_length`] = getReadLength(seq)
	flowcell[`machine_name`] = run.Instrument.Name
	flowcell[`application_name`] = run.Instrument.Type
	// flowcell[`run_param_output_folder`] = run.Href //!!! wrong information. it shall be found from RunParameters.xml
	flowcell[`description`] = run.Instrument.PlatformName
	// if len(run.DateInstrumentStarted) >= 10 {
	// 	flowcell[`run_start_date`] = run.DateInstrumentStarted[0:9]
	// }
	flowcell[`chemistry`] = run.ReagentBarcode //!!!incorrect
	flowcell[`cycles`] = getCycles(seq)
	flowcell[`current_cycle`] = seq.MaxCycleExtracted

	if t, err := dateparse.ParseLocal(run.DateInstrumentStarted); err == nil {
		if !t.IsZero() {
			flowcell[`cif_first`] = &t
			flowcell[`run_start_date`] = &t
		}

	}
	if t, err := dateparse.ParseLocal(run.DateInstrumentCompleted); err == nil {
		if !t.IsZero() {
			flowcell[`cif_latest`] = &t
		}

	}

	flowcell[`percent_pf`] = seq.PercentPf
	flowcell[`total_pf_yields_gb`] = seq.YieldTotal * seq.PercentPf
	flowcell[`mean_error_rate_r1`] = seq.ErrorRateR1
	flowcell[`mean_error_rate_r2`] = seq.ErrorRateR2
	flowcell[`mean_percent_q30`] = seq.PercentGtQ30
	flowcell[`mean_percent_aligned`] = seq.PercentAligned

	//add read summary
	reads := []map[string]interface{}{}

	for _, read := range seq.Reads {
		topushread := make(map[string]interface{})
		reads = append(reads, topushread)

		topushread[`level`] = fmt.Sprintf(`Read %d`, read.ReadNumber)
		topushread[`is_index`] = `no`
		if read.IsIndexed {
			topushread[`is_index`] = `yes`
		}
		topushread[`yield`] = read.YieldTotal
		topushread[`projected_yield`] = read.ProjectedTotalYield
		topushread[`aligned`] = read.PercentAligned
		topushread[`error_rate`] = read.ErrorRate
		topushread[`intensity_c1`] = read.IntensityCycle1
		topushread[`pct_q30`] = read.PercentGtQ30
		topushread[`pct_q30_last_10_cycles`] = read.PercentGtQ30Last10Cycles

	}
	ret[`flowcell_read_sav`] = reads

	lanereads := []map[string]interface{}{}

	isIndex := func(i int) bool {
		for _, r := range seq.Reads {
			if r.ReadNumber == i {
				return r.IsIndexed
			}
		}
		return false
	}
	for _, read := range seq.LanesByRead {
		topushread := make(map[string]interface{})
		lanereads = append(lanereads, topushread)

		topushread[`lane`] = read.LaneNumber
		topushread[`read`] = read.ReadNumber
		topushread[`surface`] = 0
		topushread[`is_fraction`] = `no`
		topushread[`is_index`] = `no`
		if isIndex(read.ReadNumber) {
			topushread[`is_index`] = `yes`
		}
		topushread[`tiles`] = read.TileCount
		topushread[`density`] = read.Density
		topushread[`cluster_pf`] = read.PercentPf
		topushread[`phasing`] = read.Phasing
		topushread[`prephasing`] = read.PrePhasing
		topushread[`phasing_equation`] = read.PhasingSlope
		topushread[`prephasing_equation`] = read.PrePhasingSlope
		topushread[`reads`] = read.Reads

		topushread[`reads_pf`] = read.ReadsPf

		topushread[`pct_q30`] = read.PercentGtQ30
		topushread[`pct_q30_last_10_cycles`] = read.PercentGtQ30Last10Cycles

		topushread[`yield`] = read.Yield
		topushread[`aligned`] = read.PercentAligned
		topushread[`error`] = read.ErrorRate
		topushread[`error_35`] = read.ErrorRate35
		topushread[`error_75`] = read.ErrorRate75
		topushread[`error_100`] = read.ErrorRate100
		topushread[`pct_occupied`] = read.Occupancy
		topushread[`intensity_c1`] = read.IntensityCycle1

		topushreadfraction := make(map[string]interface{})
		lanereads = append(lanereads, topushreadfraction)

		topushreadfraction[`lane`] = read.LaneNumber
		topushreadfraction[`read`] = read.ReadNumber
		topushreadfraction[`surface`] = 0
		topushreadfraction[`is_fraction`] = `yes`
		topushreadfraction[`is_index`] = `no`
		if isIndex(read.ReadNumber) {
			topushreadfraction[`is_index`] = `yes`
		}
		topushreadfraction[`tiles`] = 0
		topushreadfraction[`density`] = read.DensityDeviation
		topushreadfraction[`cluster_pf`] = read.PercentPfDeviation
		topushreadfraction[`phasing`] = 0
		topushreadfraction[`prephasing`] = 0
		topushreadfraction[`phasing_equation`] = read.PhasingOffset
		topushreadfraction[`prephasing_equation`] = read.PrePhasingOffset
		topushreadfraction[`reads`] = 0

		topushreadfraction[`reads_pf`] = 0

		topushreadfraction[`pct_q30`] = 0
		topushreadfraction[`pct_q30_last_10_cycles`] = 0

		topushreadfraction[`yield`] = 0
		topushreadfraction[`aligned`] = 0
		topushreadfraction[`error`] = 0
		topushreadfraction[`error_35`] = 0
		topushreadfraction[`error_75`] = 0
		topushreadfraction[`error_100`] = 0
		topushreadfraction[`pct_occupied`] = 0
		topushreadfraction[`intensity_c1`] = read.IntensityCycle1Deviation

	}
	ret[`lane_read_sav`] = lanereads
	return ret, nil
}

func GetSageSpec(ctx context.Context, client *bsshgo.Client, runId string) (map[string]interface{}, error) {
	seqstats, err := GetSequencingStats(ctx, client, runId)
	if err != nil {
		return nil, fmt.Errorf(`GetSequencingStats:%s`, err.Error())
	}
	runInfo, err := GetRunInfo(ctx, client, runId)
	if err != nil {
		return nil, fmt.Errorf(`GetRunInfo:%s`, err.Error())
	}

	if runInfoFromFile, err := GetFlowcellFromFile(ctx, client, runId); err == nil {

		return BuildSageSpec(seqstats, runInfo, runInfoFromFile)
		//dig into RunInfo.xml file

	}

	ret, err := BuildSageSpec(seqstats, runInfo, nil)
	if err != nil {
		return nil, err
	}
	//mandidate using the same convention
	//!!!very important to match up with appsession consumers assumptions
	ret[`location`] = client.GetRunHref(runId)
	return ret, nil
}
