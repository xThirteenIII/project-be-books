Ho aggiunto un mutex come campo del MeasuredWorker dal momento che Work() incrementava m.value++ ma 
l'accesso concorrente al campo non e' bloccato per la scrittura. In questo modo e' garantito l'incremento
thread safety. Nel metodo Work() quindi blocco il mutex, scrivo e poi sblocco.
Ho aggiunto anche una seconda soluzione, cambiando il tipo di value da int ad atomic.Int64.
In questo modo le operazioni di scrittura di value.Add(1) e value.Load() sono thread-safety per definizione.
Cosi il test passava, ma ci metteva comunque 20 secondi. Questo perche' in main_test.go si testava
lo SlowWorker, che aveva uno sleep di 5 secondi al suo interno. 
Ho quindi creato un worker che implementa l'interfaccia MeasuredWorker con il metodo Work() senza sleep.
In questo modo il test passa e dura tra i 0.003 e 0.004 secondi.
 
