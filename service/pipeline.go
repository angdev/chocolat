package service

import (
	"github.com/k0kubun/pp"
	"labix.org/v2/mgo"
)

type PipeStage map[string]interface{}
type Pipe []PipeStage

func NewPipe(stages ...PipeStage) Pipe {
	var pipe Pipe
	return append(pipe, stages...)
}

type PipelineResult []map[string]interface{}

func (this PipelineResult) CollapseField() PipelineResult {
	var collapsed PipelineResult
	for _, result := range this {
		pp.Println(collapseField(result))
		collapsed = append(collapsed, collapseField(result))
	}
	return collapsed
}

type Pipeline struct {
	pipes      []Pipe
	collection *mgo.Collection
}

func NewPipeline(c *mgo.Collection) *Pipeline {
	return &Pipeline{collection: c}
}

func (p *Pipeline) Copy() *Pipeline {
	var copied Pipeline = *p
	return &copied
}

func (p *Pipeline) Append(pipe Pipe) *Pipeline {
	if pipe == nil {
		return p
	}

	p.pipes = append(p.pipes, pipe)
	return p
}

func (p *Pipeline) Prepend(pipe Pipe) *Pipeline {
	if pipe == nil {
		return p
	}

	p.pipes = append([]Pipe{pipe}, p.pipes...)
	return p
}

func (p *Pipeline) Stages() []PipeStage {
	var stages []PipeStage
	for _, pipe := range p.pipes {
		stages = append(stages, pipe...)
	}
	return stages
}

func (p *Pipeline) Result() (PipelineResult, error) {
	var result PipelineResult
	if err := p.collection.Pipe(p.Stages()).All(&result); err != nil {
		return nil, err
	} else {
		return result.CollapseField(), nil
	}
}
