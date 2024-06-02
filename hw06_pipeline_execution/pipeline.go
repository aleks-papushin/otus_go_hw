package hw06pipelineexecution

import (
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// Попытки пройти тескейс "simple case"

// Стадии не выполняются по порядку, неочевидно, почему.
func ExecutePipelineBroken01(in In, done In, stages ...Stage) Out {
	out := make(Bi, 1)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer func() {
			close(out)
			wg.Done()
		}()

		out <- <-in

		for i := 0; i < len(stages); i++ {
			out <- <-stages[i](out) // неясно, почему дважды выполняется stages[0] и после этого зависает
		}
	}()

	wg.Wait()
	return out
}

// Отрефакторенный вариант ExecutePipelineBroken01. Работает так же (не работает)
func ExecutePipelineBroken02(in In, done In, stages ...Stage) Out {
	out := make(Bi, 1)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer func() {
			close(out)
			wg.Done()
		}()

		v := <-in
		out <- v

		for _, s := range stages { // аналогично ExecutePipelineBroken01, т.к. логика такая же
			v2 := <-s(out)
			out <- v2
		}
	}()

	wg.Wait()
	return out
}

// Упрощённый вариант ExecutePipelineBroken02. Результат тот же, так что видимо
// горутины в предыдущих 2х вариантах погоды не делали.
func ExecutePipelineBroken03(in In, done In, stages ...Stage) Out {
	out := make(Bi, 1)
	v := <-in
	out <- v
	wg := sync.WaitGroup{}
	wg.Add(1)

	for _, s := range stages {
		v2 := <-s(out)
		out <- v2
	}

	close(out)
	wg.Done()

	wg.Wait()
	return out
}

// Вариант от Copilot. Совершенно непонятно, почему работает,
// ведь каждому выполнению ExecutePipeline должно достаться по одному числу из канала in,
// а здесь единственный пайплайн прочитывает все значения из ch, который указывает на in по сути.
func ExecutePipelineWorking01(in In, done In, stages ...Stage) Out {
	out := make(Bi, 1)

	go func() {
		defer close(out)

		ch := in
		for _, stage := range stages {
			ch = stage(ch)
		}

		for v := range ch {
			out <- v
		}
	}()

	return out
}

// Изменил версию от Copilot, убрал заключение в горутину, тест перестал проходить.
// Что изменилось и почему код сломался - непонятно. Но на мой взгляд он и не должен был работать))
func ExecutePipelineBroken04(in In, done In, stages ...Stage) Out {
	out := make(Bi, 1)
	defer close(out)

	ch := in
	for _, stage := range stages {
		ch = stage(ch)
	}

	for v := range ch {
		out <- v
	}

	return out
}

// Упростил вариант от Copilot - убрал присвоение in другой переменной и использую саму in.
// Тест остался зелёным. Понимания, почему этот код работает, так и не появилось.
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi, 1)

	go func() {
		defer close(out)

		for _, stage := range stages {
			in = stage(in)
		}

		for v := range in {
			out <- v
		}
	}()

	return out
}
