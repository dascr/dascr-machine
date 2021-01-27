/***********************/
/*        PINS         */
/***********************/
// Ultrasonic
const int echoPin = 8;
const int trigPin = 9;
// Piezo
const int piezoPin[2] = { A0, A1 };
// Button
const int buttonPin = 2;
const int buttonLedPin = 3;

// Matrix
// Output Pins
const int PO_[4] = { 22, 24, 49, 47 };
// Input Pins
const int PI_[16] = { 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 53, 51 };

/***********************/
/*        VARIABLES    */
/***********************/
// Matrix Values
// 120 = single 20, 220 = double 20, 320 = triple 20
// 125 = single Bull, 225 = double bull
const int MATRIX_VALUES[4][16]={
  {212, 112, 209, 109, 214, 114, 211, 111, 208, 108, 000, 312, 309, 314, 311, 308},
  {216, 116, 207, 107, 219, 119, 203, 103, 217, 117, 225, 316, 307, 319, 303, 317},
  {202, 102, 215, 115, 210, 110, 206, 106, 213, 113, 125, 302, 315, 310, 306, 313},
  {204, 104, 218, 118, 201, 101, 220, 120, 205, 105, 000, 304, 318, 301, 320, 305}
};

bool bI[4][16];
bool bHitDetected = false;

// Piezo
int iValue[2] = { 0, 0 };
int iPiezoThreshold = 20;
bool bMissedDart = false;
// Button
int iButtonState = 0;
// Ultrasonic
long lDuration;
int iDistance;
int iUltrasonicThreshold = 0;
bool bUltrasonicOn = false;
bool bMotionDetected = false;
bool bUltrasonicThresholdMeasured = false;
// Input String from Raspberry Pi to Arduino
String sInputString = "";
boolean bStringComplete = false;

/***********************/
/*       Functions     */
/***********************/
int ReadUltrasonicDistance() {
  digitalWrite(trigPin, LOW);
  delayMicroseconds(2);
  digitalWrite(trigPin, HIGH);
  delayMicroseconds(10);
  digitalWrite(trigPin, LOW);
  lDuration = pulseIn(echoPin, HIGH);
  iDistance = lDuration*0.034/2;

  delay(100);
  return iDistance;
}

void CheckButton() {
  iButtonState = digitalRead(buttonPin);
  // Wenn gedrückt Ausgabe
  if (iButtonState == LOW) {
    Serial.println("b");
    delay(500);
  }
}

void EvalThrow() {
  bHitDetected = false;

  for (int x=0; x<4; x++) {
    digitalWrite(PO_[0], HIGH);
    digitalWrite(PO_[1], HIGH);
    digitalWrite(PO_[2], HIGH);
    digitalWrite(PO_[3], HIGH);
    digitalWrite(PO_[x], LOW);

    for (int y=0; y<16; y++) {
      bI[x][y]  = digitalRead(PI_[y]);
      if (bI[x][y] == 0) {
        Serial.println(MATRIX_VALUES[x][y]);
        delay(300);
        bHitDetected = true;
        // Set Bull to 0
        bI[1][10] = 1;
        bI[2][10] = 1;
      }
    }
  }
}

void CheckMissed() {
  bMissedDart = false;

  // Read both piezos
  for (int i = 0; i < 2; i++) {
    iValue[i] = analogRead(piezoPin[i]);
    iValue[i] = analogRead(piezoPin[i]);

    if (iValue[i] >= iPiezoThreshold) {
      bMissedDart = true;
    }
  }

  if (!bHitDetected && bMissedDart) {
    Serial.println("m");
    delay(300);
  }
}

void Blink(int times) {
  for (int i=0; i<times; i++) {
    digitalWrite(buttonLedPin, HIGH);
    delay(250);
    digitalWrite(buttonLedPin, LOW);
    delay(250);
  }
}

void SetUltrasonicThreshold() {
  int iTempDistance = ReadUltrasonicDistance();
  iUltrasonicThreshold = iTempDistance + 3;
  bUltrasonicThresholdMeasured = true;
}

void ProcessUltrasonic() {
  int iTempDistance = ReadUltrasonicDistance();
  if (iDistance > iUltrasonicThreshold) {
    Serial.println("u");
  }
}

void ProcessSerial() {
  if ( sInputString.indexOf("1") != -1) {
    digitalWrite(buttonLedPin, HIGH);
    bUltrasonicOn = true;
    SetUltrasonicThreshold();
  } else if (sInputString.indexOf("2") != -1) {
    digitalWrite(buttonLedPin, LOW);
    bUltrasonicOn = false;
    bMotionDetected = false;
    bUltrasonicThresholdMeasured = false;
  } else if (sInputString.indexOf("3") != -1) {
    bMotionDetected = true;
  } else if (sInputString.indexOf("4") != -1) {
    bUltrasonicThresholdMeasured = false;
  } else if (sInputString.indexOf("7") != -1) {
    Blink(7);
  }
  
  sInputString = "";
  bStringComplete = false;
}


/* Setup loop */
void setup() {
  // Pins Ultrasonic
  pinMode(trigPin, OUTPUT);
  pinMode(echoPin, INPUT);
  // Pin Button
  pinMode(buttonPin, INPUT_PULLUP);
  // Pin Button LED
  pinMode(buttonLedPin, OUTPUT);
  digitalWrite(buttonLedPin, LOW);
  // Matrix
  for (int i=0; i<4; i++) {
    pinMode(PO_[i], OUTPUT);
  }

  for (int i=0; i<16; i++) {
    pinMode(PI_[i], INPUT_PULLUP);
  }

  // Button blink 5 times
  Blink(5);
  // Start serial
  Serial.begin(9600);
}

/* Main loop */
void loop() {
  if (!bUltrasonicOn) {
    EvalThrow();
    CheckMissed();
    CheckButton();
  } else if (bUltrasonicOn) {
    if (bUltrasonicThresholdMeasured) {
      CheckButton();
      if (!bMotionDetected) {
        ProcessUltrasonic();
      } else {
        Blink(1);
      }
    } else {
      SetUltrasonicThreshold();
    }
  } else {
    Serial.println("Error in Main Loop");
  }

  if(bStringComplete)
    ProcessSerial();
}

/* Serial Events */
void serialEvent() {
  while (Serial.available()) {
    char inChar = (char)Serial.read();
    sInputString += inChar;
    if (inChar == '\n') {
      bStringComplete = true;
    }
  }
}
