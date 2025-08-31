package common

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strconv"
)

const (
    MAX_FIELDS_CSV = 5
    INDEX_CSV_FIRSTNAME = 0
    INDEX_CSV_LASTNAME = 1
    INDEX_CSV_DNI = 2
    INDEX_CSV_BIRTHDAY = 3
    INDEX_CSV_BETNUMBER = 4
)

//Processor that handles bets from csv file
type CSVBatchProcessor struct {
    file      *os.File
    reader    *csv.Reader
    batchSize int
    isEOF     bool
}


func createCSVProcessor(filepath string, batchSize int) (*CSVBatchProcessor, error){
	file, err := os.Open(filepath)
	if err != nil {
		log.Criticalf("failed to open file: %s", err)
        return nil, err
	}
	
	reader := csv.NewReader(file)

	processor := &CSVBatchProcessor{
		file: file,
		reader: reader,
		batchSize: batchSize,
		isEOF: false,
	}
	return processor, nil

}

func (p *CSVBatchProcessor) readNextBatch() ([]ClientBet, error) {
    if p.isEOF {
        return []ClientBet{}, io.EOF
    }
    
    var batch []ClientBet
    recordsRead := 0
    
    for recordsRead < p.batchSize {
        record, err := p.reader.Read()
        if err == io.EOF {
            p.isEOF = true
            log.Infof("action: reached_end_of_file | result: sucess | final_batch_size: %d", 
                len(batch))
            break
        }
        if err != nil {
            return batch, fmt.Errorf("error reading CSV record: %w", err)
        }

        bet, err := p.parseRecord(record)
        
        batch = append(batch, bet)
        recordsRead++
    }
    
    return batch, nil
}

func (p *CSVBatchProcessor) parseRecord(record []string) (ClientBet, error) {
    if len(record) < MAX_FIELDS_CSV {
        return ClientBet{}, fmt.Errorf("invalid CSV record, expected 5 fields, got %d", len(record))
    }

    dni, err := strconv.Atoi(record[INDEX_CSV_DNI])
    if err != nil {
        return ClientBet{}, fmt.Errorf("failed to parse DNI '%s': %w", record[2], err)
    }

    betNumber, err := strconv.Atoi(record[INDEX_CSV_BETNUMBER])
    if err != nil {
        return ClientBet{}, fmt.Errorf("failed to parse bet number '%s': %w", record[4], err)
    }

    bet := ClientBet{
        FirstName: record[INDEX_CSV_FIRSTNAME],
        LastName:  record[INDEX_CSV_LASTNAME],
        DNI:       dni,
        Birthday:  record[INDEX_CSV_BIRTHDAY],
        BetNumber: betNumber,
    }
    
    return bet, nil
}

func (p *CSVBatchProcessor) Close() error {
    if p.file != nil {
        log.Infof("action: csv_processor_closed | result: success")
        return p.file.Close()
    }
    return nil
}