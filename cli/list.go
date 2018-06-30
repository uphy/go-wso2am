package cli

import wso2am "github.com/uphy/go-wso2am"

// list lists paginated search result to the console.
func list(searchFunc wso2am.SearchFunc, headerFunc func(table *TableFormatter), printFunc func(entry interface{}, table *TableFormatter)) error {
	var (
		entryc = make(chan interface{})
		errc   = make(chan error)
		done   = make(chan struct{})
	)
	go func() {
		defer func() {
			close(entryc)
			close(errc)
			close(done)
		}()
		searchFunc(entryc, errc, done)
	}()

	f := newTableFormatter()
	headerFunc(f)
l:
	for {
		select {
		case entry, ok := <-entryc:
			if ok {
				printFunc(entry, f)
			} else {
				break l
			}
		case err, ok := <-errc:
			if ok {
				done <- struct{}{}
				return err
			}
			break l
		}
	}
	f.Flush()
	return nil
}
