package schema

import (
	"fmt"

	"github.com/docker/dhe-deploy/rethinkutil"
	"github.com/satori/go.uuid"
	rethink "gopkg.in/dancannon/gorethink.v2"
)

type JobLog struct {
	ID      string `gorethink:"id"`    // Randomly generated uuid for foreign reference
	Data    string `gorethink:"data"`  // (ideally) JSON encoded contents of the log
	JobID   string `gorethink:"jobID"` // ID of the job this log is from
	LineNum int    `gorethink:"lineNum"`
}

var JobLogsTable = rethinkutil.Table{
	Name:       "joblogs",
	PrimaryKey: "id",
	SecondaryIndexes: map[string][]string{
		"jobID_lineNum": {"jobID", "lineNum"}, // For quickly getting all logs for a job in order
		"lineNum":       nil,                  // For quickly sorting by line number once a job is selected
	},
}

func (m *jobrunnerManager) GetJobLogs(jobID string, offset, limit uint) ([]JobLog, error) {
	query := JobLogsTable.Term(m.DB())

	if limit > 0 {
		query = query.Between([]interface{}{jobID, offset}, []interface{}{jobID, offset + limit}, rethink.BetweenOpts{Index: "jobID_lineNum"})
	} else {
		query = query.Between([]interface{}{jobID, offset}, []interface{}{jobID, rethink.MaxVal}, rethink.BetweenOpts{Index: "jobID_lineNum"})
	}

	cursor, err := query.OrderBy("lineNum").Run(m.Session())
	if err != nil {
		return nil, fmt.Errorf("unable to query db for job logs: %s", err)
	}
	jobLogs := []JobLog{}
	err = cursor.All(&jobLogs)
	if err != nil {
		return nil, fmt.Errorf("unable to scan job logs: %s", err)
	}
	return jobLogs, nil
}

func (m *jobrunnerManager) InsertJobLog(jobID string, log string, lineNum int) error {
	jobLog := &JobLog{
		ID:      uuid.NewV4().String(),
		JobID:   jobID,
		Data:    log,
		LineNum: lineNum,
	}
	_, err := JobLogsTable.Term(m.DB()).Insert(jobLog).RunWrite(m.Session())
	if err != nil {
		return fmt.Errorf("unable to insert job log into db: %s", err)
	}
	return nil
}
