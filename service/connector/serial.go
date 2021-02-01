package connector

// processSerial will handle the output coming from serial
func (c *Service) startSerial() {
	go c.Read()
	for {
		select {
		case <-c.Quit:
			return
		}
	}
}

// Write will write to the serial connection
func (c *Service) Write(input string) {
	b := []byte("<" + input + ">")

	_, err := c.Conn.Write(b)
	if err != nil {
		return
	}
}

// Read will read from the serial connection
func (c *Service) Read() {
	for c.Scanner.Scan() {
		cmd := c.Scanner.Text()
		if cmd != "" {
			c.Command <- cmd
		}
	}
}
