package pipeline

import "bytes"

func split(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	var (
		ds  int
		dlt = dropStartCRLF(data)
		dll = len(data) - len(dlt)
	)

	dataRst := data

	for {
		if de := bytes.IndexByte(dataRst, '\n'); de >= 0 {
			if de >= 1 {
				di := ds + de
				d := data[ds:di]
				drt := dropEndCR(d)

				drl := len(d) - len(drt)

				switch {
				case de > 1 && dlt[0] == '-' && dlt[1] == '-': // Start with '--'.
					return di + 1, data[dll : di-drl], nil
				case de > 1 && dlt[0] == '/' && dlt[1] == '*': // Start with '/*'.
					// End with '*/'.
					if drt[de-drl-2] == '*' && drt[de-drl-1] == '/' {
						return di + 1, data[dll : di-drl], nil
					}

					// Revert searching.
					var ending bool

					for i := len(drt) - 1; i >= 0; i-- {
						if drt[i] == '/' {
							ending = true
							continue
						}

						if ending {
							if drt[i] == '*' {
								// Comment command.
								return di + 1, data[dll : di-drl], nil
							}

							break
						}
					}
				case drt[de-drl-1] == ';': // End with ';'.
					return di + 1, data[dll : di-drl], nil
				}
			}

			if de+1 >= len(dataRst) {
				break
			}
			dataRst = dataRst[de+1:]
			ds += de + 1

			continue
		}

		break
	}

	if atEOF {
		return len(data), dropEndCR(data[dll:]), nil
	}

	return 0, nil, nil
}

// dropStartCRLF drops leading \r\n from the data.
func dropStartCRLF(data []byte) []byte {
	return bytes.TrimLeft(data, "\r\n")
}

// dropEndCR drops terminal \r from the data.
func dropEndCR(data []byte) []byte {
	return bytes.TrimRight(data, "\r")
}
