package concurrency

/*
Processor is a struct meant to queue up values that should be processed with a supplied processFunc
*/
type Processor struct {
	process     func(val interface{})
	processChan chan interface{}
	stopChan    chan bool
	running     bool
	Stopped     bool
	DoneChan    chan bool
}

/*
NewProcessor returns a new processor that will start upon the first call to Push.
*/
func NewProcessor(processFunc func(val interface{})) *Processor {
	return &Processor{
		process:     processFunc,
		processChan: make(chan interface{}),
		stopChan:    make(chan bool),
	}
}

/*
Push is a non-blocking call that will queue up another value to be processed.
If the processor is not yet running, it will be started.
*/
func (processor *Processor) Push(val interface{}) {
	if !processor.running {
		processor.running = true
		go processor.run()
	}
	go processor.sendToProcessChan(val)
}

/*
Stop will stop the processor. IMPORTANT NOTE: Stop will cut off the processor before all values are finished being processed.
Have the last call to process trigger this function.
*/
func (processor *Processor) Stop() {
	processor.stopChan <- true
}

func (processor *Processor) sendToProcessChan(val interface{}) {
	processor.processChan <- val
}

func (processor *Processor) run() {
	processor.running = true
	keepRunning := true
	for keepRunning {
		select {
		case <-processor.stopChan:
			keepRunning = false
			processor.Stopped = true
			if processor.DoneChan != nil {
				processor.DoneChan <- true
			}
		case val := <-processor.processChan:
			processor.process(val)
		}
	}
}
