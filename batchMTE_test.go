// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package ach

import (
	"log"
	"testing"
	"time"

	"github.com/moov-io/base"
)

// mockBatchMTEHeader creates a MTE batch header
func mockBatchMTEHeader() *BatchHeader {
	bh := NewBatchHeader()
	bh.ServiceClassCode = DebitsOnly
	bh.CompanyName = "Merchant with ATM"
	bh.CompanyIdentification = "231380104"
	bh.StandardEntryClassCode = MTE
	bh.CompanyEntryDescription = "CASH WITHDRAW"
	bh.EffectiveEntryDate = base.NewTime(time.Now())
	bh.ODFIIdentification = "23138010"
	return bh
}

// mockMTEEntryDetail creates a MTE entry detail
func mockMTEEntryDetail() *EntryDetail {
	entry := NewEntryDetail()
	entry.TransactionCode = CheckingDebit
	entry.SetRDFI("031300012")
	entry.DFIAccountNumber = "744-5678-99"
	entry.Amount = 10000
	entry.SetOriginalTraceNumber("031300010000001")
	entry.SetReceivingCompany("JANE DOE")
	entry.SetTraceNumber("23138010", 1)
	entry.AddendaRecordIndicator = 1

	addenda02 := NewAddenda02()

	// NACHA rules example: 200509*321 East Market Street*Anytown*VA\
	addenda02.TerminalIdentificationCode = "200509"
	addenda02.TerminalLocation = "321 East Market Street"
	addenda02.TerminalCity = "ANYTOWN"
	addenda02.TerminalState = "VA"

	addenda02.TransactionSerialNumber = "123456" // Generated by Terminal, used for audits
	addenda02.TransactionDate = "1224"
	addenda02.TraceNumber = entry.TraceNumber
	entry.Addenda02 = addenda02

	return entry
}

// mockBatchMTE creates a MTE batch
func mockBatchMTE() *BatchMTE {
	batch := NewBatchMTE(mockBatchMTEHeader())
	batch.AddEntry(mockMTEEntryDetail())
	if err := batch.Create(); err != nil {
		log.Fatalf("Unexpected error building batch: %s\n", err)
	}
	return batch
}

// testBatchMTEHeader creates a MTE batch header
func testBatchMTEHeader(t testing.TB) {
	batch, _ := NewBatch(mockBatchMTEHeader())
	_, ok := batch.(*BatchMTE)
	if !ok {
		t.Error("Expecting BatchMTE")
	}
}

// TestBatchMTEHeader tests creating a MTE batch header
func TestBatchMTEHeader(t *testing.T) {
	testBatchMTEHeader(t)
}

// BenchmarkBatchMTEHeader benchmark creating a MTE batch header
func BenchmarkBatchMTEHeader(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBatchMTEHeader(b)
	}
}

// testBatchMTEAddendumCount batch control MTE can only have one addendum per entry detail
func testBatchMTEAddendumCount(t testing.TB) {
	mockBatch := mockBatchMTE()
	// Adding a second addenda to the mock entry
	mockBatch.GetEntries()[0].Addenda02 = mockAddenda02()
	if err := mockBatch.Validate(); err != nil {
		if e, ok := err.(*BatchError); ok {
			if e.FieldName != "EntryAddendaCount" {
				t.Errorf("%T: %s", err, err)
			}
		} else {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestBatchMTEAddendumCount tests batch control MTE can only have one addendum per entry detail
func TestBatchMTEAddendumCount(t *testing.T) {
	testBatchMTEAddendumCount(t)
}

// BenchmarkBatchMTEAddendumCount benchmarks batch control MTE can only have one addendum per entry detail
func BenchmarkBatchMTEAddendumCount(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBatchMTEAddendumCount(b)
	}
}

// TestBatchMTEAddendum02 validates Addenda02 returns an error
func TestBatchMTEAddendum02(t *testing.T) {
	mockBatch := NewBatchMTE(mockBatchMTEHeader())
	mockBatch.AddEntry(mockMTEEntryDetail())
	if err := mockBatch.Create(); err != nil {
		if e, ok := err.(*BatchError); ok {
			if e.FieldName != "TypeCode" {
				t.Errorf("%T: %s", err, err)
			}
		} else {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// testBatchMTEReceivingCompanyName validates Receiving company / Individual name is a mandatory field
func testBatchMTEReceivingCompanyName(t testing.TB) {
	mockBatch := mockBatchMTE()
	// modify the Individual name / receiving company to nothing
	mockBatch.GetEntries()[0].SetReceivingCompany("")
	if err := mockBatch.Validate(); err != nil {
		if e, ok := err.(*BatchError); ok {
			if e.FieldName != "IndividualName" {
				t.Errorf("%T: %s", err, err)
			}
		} else {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestBatchMTEReceivingCompanyName tests validating receiving company / Individual name is a mandatory field
func TestBatchMTEReceivingCompanyName(t *testing.T) {
	testBatchMTEReceivingCompanyName(t)
}

// BenchmarkBatchMTEReceivingCompanyName benchmarks validating receiving company / Individual name is a mandatory field
func BenchmarkBatchMTEReceivingCompanyName(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBatchMTEReceivingCompanyName(b)
	}
}

// testBatchMTEAddendaTypeCode validates addenda type code is 05
func testBatchMTEAddendaTypeCode(t testing.TB) {
	mockBatch := mockBatchMTE()
	mockBatch.GetEntries()[0].Addenda02.TypeCode = "05"
	if err := mockBatch.Validate(); err != nil {
		if e, ok := err.(*BatchError); ok {
			if e.FieldName != "TypeCode" {
				t.Errorf("%T: %s", err, err)
			}
		} else {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestBatchMTEAddendaTypeCode tests validating addenda type code is 05
func TestBatchMTEAddendaTypeCode(t *testing.T) {
	testBatchMTEAddendaTypeCode(t)
}

// BenchmarkBatchMTEAddendaTypeCod benchmarks validating addenda type code is 05
func BenchmarkBatchMTEAddendaTypeCode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBatchMTEAddendaTypeCode(b)
	}
}

// testBatchMTESEC validates that the standard entry class code is MTE for batchMTE
func testBatchMTESEC(t testing.TB) {
	mockBatch := mockBatchMTE()
	mockBatch.Header.StandardEntryClassCode = ACK
	if err := mockBatch.Validate(); err != nil {
		if e, ok := err.(*BatchError); ok {
			if e.FieldName != "StandardEntryClassCode" {
				t.Errorf("%T: %s", err, err)
			}
		} else {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestBatchMTESEC tests validating that the standard entry class code is MTE for batchMTE
func TestBatchMTESEC(t *testing.T) {
	testBatchMTESEC(t)
}

// BenchmarkBatchMTESEC benchmarks validating that the standard entry class code is MTE for batch MTE
func BenchmarkBatchMTESEC(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBatchMTESEC(b)
	}
}

// testBatchMTEServiceClassCode validates ServiceClassCode
func testBatchMTEServiceClassCode(t testing.TB) {
	mockBatch := mockBatchMTE()
	// Batch Header information is required to Create a batch.
	mockBatch.GetHeader().ServiceClassCode = 0
	mockBatch.Create()
	if err := mockBatch.Validate(); err != nil {
		if e, ok := err.(*BatchError); ok {
			if e.FieldName != "ServiceClassCode" {
				t.Errorf("%T: %s", err, err)
			}
		} else {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestBatchMTEServiceClassCode tests validating ServiceClassCode
func TestBatchMTEServiceClassCode(t *testing.T) {
	testBatchMTEServiceClassCode(t)
}

// BenchmarkBatchMTEServiceClassCode benchmarks validating ServiceClassCode
func BenchmarkBatchMTEServiceClassCode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBatchMTEServiceClassCode(b)
	}
}

// TestBatchMTEAmount validates Amount
func TestBatchMTEAmount(t *testing.T) {
	mockBatch := mockBatchMTE()
	mockBatch.GetEntries()[0].Amount = 0
	mockBatch.Create()
	if err := mockBatch.Validate(); err != nil {
		if e, ok := err.(*BatchError); ok {
			if e.FieldName != "Amount" {
				t.Errorf("%T: %s", err, err)
			}
		} else {
			t.Errorf("%T: %s", err, err)
		}
	}
}

func TestBatchMTETerminaalState(t *testing.T) {
	mockBatch := mockBatchMTE()
	mockBatch.GetEntries()[0].Addenda02.TerminalState = "XX"
	mockBatch.Create()
	if err := mockBatch.Validate(); err != nil {
		if e, ok := err.(*BatchError); ok {
			if e.FieldName != "TerminalState" {
				t.Errorf("%T: %s", err, err)
			}
		} else {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestBatchMTEIndividualName validates IndividualName
func TestBatchMTEIndividualName(t *testing.T) {
	mockBatch := mockBatchMTE()
	mockBatch.GetEntries()[0].IndividualName = ""
	mockBatch.Create()
	if err := mockBatch.Validate(); err != nil {
		if e, ok := err.(*BatchError); ok {
			if e.FieldName != "IndividualName" {
				t.Errorf("%T: %s", err, err)
			}
		} else {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestBatchMTEIdentificationNumber validates IdentificationNumber
func TestBatchMTEIdentificationNumber(t *testing.T) {
	mockBatch := mockBatchMTE()

	// NACHA rules state MTE records can't be all spaces or all zeros
	mockBatch.GetEntries()[0].IdentificationNumber = "   "
	mockBatch.Create()
	if err := mockBatch.Validate(); err != nil {
		if e, ok := err.(*BatchError); ok {
			if e.FieldName != "IdentificationNumber" {
				t.Errorf("%T: %s", err, err)
			}
		} else {
			t.Errorf("%T: %s", err, err)
		}
	}

	mockBatch.GetEntries()[0].IdentificationNumber = "000000"
	mockBatch.Create()
	if err := mockBatch.Validate(); err != nil {
		if e, ok := err.(*BatchError); ok {
			if e.FieldName != "IdentificationNumber" {
				t.Errorf("%T: %s", err, err)
			}
		} else {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestBatchMTEValidTranCodeForServiceClassCode validates a transactionCode based on ServiceClassCode
func TestBatchMTEValidTranCodeForServiceClassCode(t *testing.T) {
	mockBatch := mockBatchMTE()
	mockBatch.GetHeader().ServiceClassCode = CreditsOnly
	if err := mockBatch.Create(); err != nil {
		if e, ok := err.(*BatchError); ok {
			if e.FieldName != "TransactionCode" {
				t.Errorf("%T: %s", err, err)
			}
		} else {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestBatchMTEAddenda05 validates BatchMTE cannot have Addenda05
func TestBatchMTEAddenda05(t *testing.T) {
	mockBatch := mockBatchMTE()
	mockBatch.Entries[0].AddendaRecordIndicator = 1
	mockBatch.GetEntries()[0].AddAddenda05(mockAddenda05())
	mockBatch.Create()
	if err := mockBatch.Validate(); err != nil {
		if e, ok := err.(*BatchError); ok {
			if e.FieldName != "Addenda05" {
				t.Errorf("%T: %s", err, err)
			}
		} else {
			t.Errorf("%T: %s", err, err)
		}
	}
}
