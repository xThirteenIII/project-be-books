# Solution
Ho reso MeasuredWorker thread-safe aggiungendo un mutex sul campo value, così Work() può incrementarlo in modo sicuro con accesso esclusivo.
In alternativa, ho provato anche a sostituire int con atomic.Int64, usando Add(1) e Load() per garantire la sicurezza concorrente senza mutex.
Nel test ho poi evitato di usare SlowWorker, perché il suo sleep di 5 secondi faceva durare l’esecuzione circa 20 secondi.
Al suo posto ho creato workerFunc, un tipo funzione che implementa Work() e soddisfa l’interfaccia Worker, così nei test posso passare un worker finto che non fa nulla ma permette di testare solo la logica di conteggio.
In questo modo il test resta corretto e scende a circa 3–4 millisecondi.

# Golang concurrent test

We're writing a `MeasuredWorker` that tries to parallelize and count the work done by a generic worker defined by the
interface `Worker`.

A test has been written, however we realized two problems:

- when processing concurrently the count not always matches the expected number of operations done
- the test requires over 20 seconds to run

## Constraints

You're allowed to change only the `measured_worker.go` and the `main_test.go` files.

## Instructions

If you have Go installed you can run the tests by simply doing:

```bash
go test .
```

Otherwise, you can use the included Dockerfile to build and run the tests in a container:

```bash
docker run --rm -it $(docker build -q .)
```

## Tips

Remember that the focus is `MeasuredWorker`, so feel free to mock or workaround the other components that are in the
files you're not allowed to edit!
