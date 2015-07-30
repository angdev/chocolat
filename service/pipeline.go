package service

type PipeStage map[string]interface{}
type Pipe []PipeStage

func NewPipe(stages ...PipeStage) Pipe {
	var pipe Pipe
	return append(pipe, stages...)
}

func (this Pipe) Join(pipes ...Pipe) Pipe {
	joined := this
	for _, pipe := range pipes {
		joined = append(joined, pipe...)
	}
	return joined
}

type PipelineResult []map[string]interface{}

func (this PipelineResult) CollapseField() PipelineResult {
	var collapsed PipelineResult
	for _, result := range this {
		collapsed = append(collapsed, collapseField(result))
	}
	return collapsed
}

type Pipeline struct {
	pipes []Pipe
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
