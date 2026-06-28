package main

import "sync"

// SOLUZIONE 1

type MeasuredWorker struct {
	Worker
	value int
	mu    sync.Mutex
}

func (m *MeasuredWorker) Work() {
	m.Worker.Work()
	m.mu.Lock()
	defer m.mu.Unlock()

	m.value++
}

func (m *MeasuredWorker) Value() int {

	return m.value
}

/* SOLUZIONE 2

type MeasuredWorker struct {
	Worker
	value atomic.Int64
}

func (m *MeasuredWorker) Work() {
	m.Worker.Work()

	m.value.Add(1)
}

func (m *MeasuredWorker) Value() int {

	return int(m.value.Load())
}
*/
