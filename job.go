package main

import (
	"bytes"
	"encoding/json"
	"strconv"
)

type Job struct {
	Raw       *jobResponse
	Requester *Requester
	Base      string
}

type Cause struct {
	ShortDescription string
	UserId           string
	Username         string
}

type ActionsObject struct {
	FailCount  int64
	SkipCount  int64
	TotalCount int64
	UrlName    string
}

type jobBuild struct {
	Number int
	Url    string
}

type jobResponse struct {
	Actions            interface{}
	Buildable          bool `json:"buildable"`
	Builds             []jobBuild
	Color              string        `json:"color"`
	ConcurrentBuild    bool          `json:"concurrentBuild"`
	Description        string        `json:"description"`
	DisplayName        string        `json:"displayName"`
	DisplayNameOrNull  interface{}   `json:"displayNameOrNull"`
	DownstreamProjects []interface{} `json:"downstreamProjects"`
	FirstBuild         jobBuild
	HealthReport       []struct {
		Description   string  `json:"description"`
		IconClassName string  `json:"iconClassName"`
		IconUrl       string  `json:"iconUrl"`
		Score         float64 `json:"score"`
	} `json:"healthReport"`
	InQueue               bool     `json:"inQueue"`
	KeepDependencies      bool     `json:"keepDependencies"`
	LastBuild             jobBuild `json:"lastBuild"`
	LastCompletedBuild    jobBuild `json:"lastCompletedBuild"`
	LastFailedBuild       jobBuild `json:"lastFailedBuild"`
	LastStableBuild       jobBuild `json:"lastStableBuild"`
	LastSuccessfulBuild   jobBuild `json:"lastSuccessfulBuild"`
	LastUnstableBuild     jobBuild `json:"lastUnstableBuild"`
	LastUnsuccessfulBuild jobBuild `json:"lastUnsuccessfulBuild"`
	Name                  string   `json:"name"`
	NextBuildNumber       float64  `json:"nextBuildNumber"`
	Property              []struct {
		ParameterDefinitions []struct {
			DefaultParameterValue struct {
				Name  string `json:"name"`
				Value bool   `json:"value"`
			} `json:"defaultParameterValue"`
			Description string `json:"description"`
			Name        string `json:"name"`
			Type        string `json:"type"`
		} `json:"parameterDefinitions"`
	} `json:"property"`
	QueueItem        interface{}   `json:"queueItem"`
	Scm              struct{}      `json:"scm"`
	UpstreamProjects []interface{} `json:"upstreamProjects"`
	URL              string        `json:"url"`
}

func (j *Job) GetName() string {
	return j.Raw.Name
}

func (j *Job) GetDescription() string {
	return j.Raw.Description
}

func (j *Job) GetDetails() *jobResponse {
	return j.Raw
}

func (j *Job) GetBuild(id string) *Build {
	build := Build{Raw: new(buildResponse), Requester: j.Requester, Base: "/job/" + j.GetName() + "/" + number}
	if build.Poll() == 200 {
		return &build
	}
	return nil
}

func (j *Job) getBuildByType(buildType string) *Build {
	allowed := map[string]jobBuild{
		"lastStableBuild":     j.Raw.LastStableBuild,
		"lastSuccessfulBuild": j.Raw.LastSuccessfulBuild,
		"lastBuild":           j.Raw.LastBuild,
		"lastCompletedBuild":  j.Raw.LastCompletedBuild,
		"firstBuild":          j.Raw.FirstBuild,
		"lastFailedBuild":     j.Raw.LastFailedBuild,
	}
	number := ""
	if val, ok := allowed[buildType]; ok {
		number = strconv.Itoa(val.Number)
	} else {
		panic("No Such Build")
	}
	build := Build{Raw: new(buildResponse), Requester: j.Requester, Base: "/job/" + j.GetName() + "/" + number}
	if build.Poll() == 200 {
		return &build
	}
	return nil
}

func (j *Job) GetLastSuccessfulBuild() *Build {
	return j.getBuildByType("lastSuccessfulBuild")
}

func (j *Job) GetFirstBuild() *Build {
	return j.getBuildByType("firstBuild")
}

func (j *Job) GetLastBuild() *Build {
	return j.getBuildByType("lastBuild")
}

func (j *Job) GetLastStableBuild() *Build {
	return j.getBuildByType("lastStableBuild")
}

func (j *Job) GetLastFailedBuild() *Build {
	return j.getBuildByType("lastFailedBuild")
}

func (j *Job) GetLastCompletedBuild() *Build {
	return j.getBuildByType("lastCompletedBuild")
}

func (j *Job) GetAllBuilds() {
	j.Poll()
	builds := make([]*Build, len(j.Raw.Builds))
	for q, v := range j.Raw.Builds {

	}
}

func (j *Job) GetBuildMetaData() {

}

func (j *Job) GetUpstreamJobNames() {

}

func (j *Job) GetDownstreamJobNames() {

}

func (j *Job) GetUpstreamJobs() {

}

func (J *Job) GetDownstreamJobs() {

}

func (j *Job) Enable() bool {
	resp := j.Requester.Post(j.Base+"/enable", nil, nil, nil)
	return resp.StatusCode == 200
}

func (j *Job) Disable() bool {
	resp := j.Requester.Post(j.Base+"/disable", nil, nil, nil)
	return resp.StatusCode == 200
}

func (j *Job) Delete() bool {
	resp := j.Requester.Post(j.Base+"/doDelete", nil, nil, nil)
	return resp.StatusCode == 200
}

func (j *Job) Rename(name string) {
	payload, _ := json.Marshal(map[string]string{"newName": name})
	j.Requester.Post(j.Base+"/doRename", bytes.NewBuffer(payload), nil, nil)
}

func (j *Job) Exists() {

}

func (j *Job) Create(config string) *Job {
	resp := j.Requester.Post("/createItem", bytes.NewBuffer([]byte(config)), j.Raw, nil)
	if resp.Status == "200" {
		return j
	} else {
		return nil
	}
}

func (j *Job) Copy(from string, newName string) *Job {
	qr := map[string]string{"name": newName, "from": from, "mode": "copy"}
	resp := j.Requester.Post("/createItem", nil, nil, qr)
	if resp.StatusCode == 200 {
		return j
	}
	return nil
}

func (j *Job) GetConfig() string {
	var data string
	j.Requester.GetXML(j.Base+"/config.xml", &data, nil)
	return data
}

func (j *Job) SetConfig() {

}

func (j *Job) IsQueued() bool {
	j.Poll()
	return j.Raw.InQueue
}

func (j *Job) IsRunning() {
	j.Poll()
}

func (j *Job) IsEnabled() bool {
	j.Poll()
	return j.Raw.Color != "disabled"
}

func (j *Job) HasQueuedBuild() {

}

func (j *Job) Invoke(files []string, options ...interface{}) bool {
	return true
}

func (j *Job) Poll() int {
	j.Requester.Get(j.Base, j.Raw, nil)
	return j.Requester.LastResponse.StatusCode
}
