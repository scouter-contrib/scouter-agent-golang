package util

import (
	"testing"
)

func TestStrMatch(t *testing.T) {

	anyAsta := NewStrMatch("*")
	startAsta := NewStrMatch("*HelloWorld")
	endAsta := NewStrMatch("HelloWorld*")
	midAsta := NewStrMatch("Hello*World")
	startMidAsta := NewStrMatch("*Hello*World")
	midEndAsta := NewStrMatch("Hello*World*")
	startEndAsta := NewStrMatch("*HelloWorld*")
	complexAsta := NewStrMatch("*Hello*World*")
	noAasta := NewStrMatch("HelloWorld")

	var testAny = "Xxx"
	assertTrue(t, anyAsta.include(testAny))
	assertFalse(t, startAsta.include(testAny))
	assertFalse(t, endAsta.include(testAny))
	assertFalse(t, midAsta.include(testAny))
	assertFalse(t, startMidAsta.include(testAny))
	assertFalse(t, midEndAsta.include(testAny))
	assertFalse(t, startEndAsta.include(testAny))
	assertFalse(t, complexAsta.include(testAny))
	assertFalse(t, noAasta.include(testAny))

	var testExact = "HelloWorld"
	assertTrue(t, anyAsta.include(testExact))
	assertTrue(t, startAsta.include(testExact))
	assertTrue(t, endAsta.include(testExact))
	assertTrue(t, midAsta.include(testExact))
	assertTrue(t, startMidAsta.include(testExact))
	assertTrue(t, midEndAsta.include(testExact))
	assertTrue(t, startEndAsta.include(testExact))
	assertTrue(t, complexAsta.include(testExact))
	assertTrue(t, noAasta.include(testExact))

	var testStart = "HelloWorldxxx"
	assertTrue(t, anyAsta.include(testStart))
	assertFalse(t, startAsta.include(testStart))
	assertTrue(t, endAsta.include(testStart))
	assertFalse(t, midAsta.include(testStart))
	assertFalse(t, startMidAsta.include(testStart))
	assertTrue(t, midEndAsta.include(testStart))
	assertTrue(t, startEndAsta.include(testStart))
	assertTrue(t, complexAsta.include(testStart))
	assertFalse(t, noAasta.include(testStart))

	var testEnd = "XxxHelloWorld"
	assertTrue(t, anyAsta.include(testEnd))
	assertTrue(t, startAsta.include(testEnd))
	assertFalse(t, endAsta.include(testEnd))
	assertFalse(t, midAsta.include(testEnd))
	assertTrue(t, startMidAsta.include(testEnd))
	assertFalse(t, midEndAsta.include(testEnd))
	assertTrue(t, startEndAsta.include(testEnd))
	assertTrue(t, complexAsta.include(testEnd))
	assertFalse(t, noAasta.include(testEnd))

	var testMid = "XxxHelloWorldXxx"
	assertTrue(t, anyAsta.include(testMid))
	assertFalse(t, startAsta.include(testMid))
	assertFalse(t, endAsta.include(testMid))
	assertFalse(t, midAsta.include(testMid))
	assertFalse(t, startMidAsta.include(testMid))
	assertFalse(t, midEndAsta.include(testMid))
	assertTrue(t, startEndAsta.include(testMid))
	assertTrue(t, complexAsta.include(testMid))
	assertFalse(t, noAasta.include(testMid))

	var testStartMid = "HelloXxxWorldxxx"
	assertTrue(t, anyAsta.include(testStartMid))
	assertFalse(t, startAsta.include(testStartMid))
	assertFalse(t, endAsta.include(testStartMid))
	assertFalse(t, midAsta.include(testStartMid))
	assertFalse(t, startMidAsta.include(testStartMid))
	assertTrue(t, midEndAsta.include(testStartMid))
	assertFalse(t, startEndAsta.include(testStartMid))
	assertTrue(t, complexAsta.include(testStartMid))
	assertFalse(t, noAasta.include(testStartMid))

	var testMidMid = "xxxHelloXxxWorldxxx"
	assertTrue(t, anyAsta.include(testMidMid))
	assertFalse(t, startAsta.include(testMidMid))
	assertFalse(t, endAsta.include(testMidMid))
	assertFalse(t, midAsta.include(testMidMid))
	assertFalse(t, startMidAsta.include(testMidMid))
	assertFalse(t, midEndAsta.include(testMidMid))
	assertFalse(t, startEndAsta.include(testMidMid))
	assertTrue(t, complexAsta.include(testMidMid))
	assertFalse(t, noAasta.include(testMidMid))

	var testMidEnd = "xxxHelloXxxWorld"
	assertTrue(t, anyAsta.include(testMidEnd))
	assertFalse(t, startAsta.include(testMidEnd))
	assertFalse(t, endAsta.include(testMidEnd))
	assertFalse(t, midAsta.include(testMidEnd))
	assertTrue(t, startMidAsta.include(testMidEnd))
	assertFalse(t, midEndAsta.include(testMidEnd))
	assertFalse(t, startEndAsta.include(testMidEnd))
	assertTrue(t, complexAsta.include(testMidEnd))
	assertFalse(t, noAasta.include(testMidEnd))

}

func assertTrue(t *testing.T, b bool) {
	if !b {
		//t.Error("match is not true! ")
		//t.Log(string(debug.Stack()))
		panic("match is not true! ")
	}
}

func assertFalse(t *testing.T, b bool) {
	if b {
		panic("match is not false! ")
	}
}
