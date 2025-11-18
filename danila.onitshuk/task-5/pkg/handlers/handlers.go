package handlers

import (
	"context"
	"strings"
)

/*
Модификатор данных, выполняющий добавление к входных данным
префикса “decorated: <исходные данные>”, если данный префикс не
был добавлен ранее. Если входные данные содержат подстроку “no
decorator” – обработчик должен завершить выполнение и вернуть
ошибку, содержащую подстроку “can’t be decorated”.
*/
func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case line, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(line, noDecoratorData) {
				return ErrNoDecorator
			} else if !strings.Contains(line, textForDecoratorString) {
				line = textForDecoratorString + line
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case output <- line:
			}
		}
	}
}

/*
Мультиплексор, который принимает на вход данные и объединяет их.
В случае, если данные содержат подстроку “no multiplexer” – данные
должны быть пропущены (фильтрация данных).
*/
func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {
	defer close(output)

	for _, inputChan := range inputs {
		go func(ch <-chan string) {
			for {
				select {
				case <-ctx.Done():
					return
				case line, ok := <-ch:
					if !ok {
						return
					}
					if !strings.Contains(line, noMultiplexerData) {
						select {
						case <-ctx.Done():
							return
						case output <- line:
						}
					}
				}
			}
		}(inputChan)
	}

	<-ctx.Done()
	return ctx.Err()
}

/*
Сепаратор, выполняющий разделение данных по каналам на основе
порядкового номера их получения. Например, для двух выходных
каналов: первое значение будет передано в первый выходной канал,
второе – во второй, третье – в первый и тд.
*/
func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	defer (func() {
		for _, ch := range outputs {
			close(ch)
		}
	})()

	var (
		cnt    = 0
		cntOut = len(outputs)
	)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case line, ok := <-input:
			if !ok {
				return nil
			}

			if cntOut != 0 {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case outputs[cnt%cntOut] <- line:
					cnt++
				}
			}
		}
	}
}
