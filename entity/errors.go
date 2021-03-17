package entity

import "errors"

var ErrCreateChecksumNotSuccesfull = errors.New("Creation of checksum not succesful")
var ErrOpeningOfFileUnsuccesfull = errors.New("Opening of file was unscuccsesful")
var ErrChildAlreadyExists = errors.New("Child already exists")
var ErrUnsuccesfulValidationOfDirectory = errors.New("Validation fo client failed")
var ErrFileNotFound = errors.New("Could not find file")
var ErrFailedDirectoryScan = errors.New("Failed the reading of the directory")
var ErrCouldNotWriteToFile = errors.New("Could not write to file")
var ErrCouldNotMarshallJSON = errors.New("Could not marshall json")
var ErrCouldNotUnMarshallJSON = errors.New("Could not unmarshall json")
var ErrFieldWasEmpty = errors.New("Field was empty")
var ErrWrongStatus = errors.New("Wrong status")
