package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		out := make(Bi)
		go func(s Stage, in In, out Bi) {
			defer close(out)
			for v := range s(in) {
				select {
				case <-done:
					return
				default:
					out <- v
				}
			}
		}(stage, in, out)
		in = out
	}
	return in
}
