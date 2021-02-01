/***********************/
/*        PINS         */
/***********************/
// Ultrasonic
const int echoPin = 8;
const int trigPin = 9;
// Piezo
const int piezoPin[2] = {A0, A1};
// Button
const int buttonPin = 2;
const int buttonLedPin = 3;

// Matrix
// Output Pins
const int PO_[4] = {22, 24, 49, 47};
// Input Pins
const int PI_[16] = {26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 53, 51};

/***********************/
/*        VARIABLES    */
/***********************/
// Matrix Values
// 120 = single 20, 220 = double 20, 320 = triple 20
// 125 = single Bull, 225 = double bull
const int MATRIX_VALUES[4][16] = {
    {212, 112, 209, 109, 214, 114, 211, 111, 208, 108, 000, 312, 309, 314, 311, 308},
    {216, 116, 207, 107, 219, 119, 203, 103, 217, 117, 225, 316, 307, 319, 303, 317},
    {202, 102, 215, 115, 210, 110, 206, 106, 213, 113, 125, 302, 315, 310, 306, 313},
    {204, 104, 218, 118, 201, 101, 220, 120, 205, 105, 000, 304, 318, 301, 320, 305}};

bool bI[4][16];
bool bHitDetected = false;

// Piezo
int iValue[2] = {0, 0};
int iPiezoThreshold = 20;
bool bMissedDart = false;
// Button
int iButtonState = 0;
// Ultrasonic
long lDuration;
int iDistance;
int iUltrasonicThreshold = 0;
bool bMotionDetected = false;
bool bUltrasonicThresholdMeasured = false;
// Input String from Raspberry Pi to Arduino
String sInputString = "";
boolean bStringComplete = false;

// State
// 0: INIT
// 1: THROW
// 2: NEXTPLAYER
// 3: MOTION DETECTED
// 4: RESET ULTRASONIC
// 5: WON
// 9: CONNECTED
int iState = 0;

/***********************/
/*       Functions     */
/***********************/
int ReadUltrasonicDistance()
{
  digitalWrite(trigPin, LOW);
  delayMicroseconds(2);
  digitalWrite(trigPin, HIGH);
  delayMicroseconds(10);
  digitalWrite(trigPin, LOW);
  lDuration = pulseIn(echoPin, HIGH);
  iDistance = lDuration * 0.034 / 2;

  delay(100);
  return iDistance;
}

void CheckButton()
{
  iButtonState = digitalRead(buttonPin);
  // Wenn gedrückt Ausgabe
  if (iButtonState == LOW)
  {
    Serial.println("b");
    delay(500);
  }
}

void ButtonOn()
{
  digitalWrite(buttonLedPin, HIGH);
}

void ButtonOff()
{
  digitalWrite(buttonLedPin, LOW);
}

void EvalThrow()
{
  bHitDetected = false;

  for (int x = 0; x < 4; x++)
  {
    digitalWrite(PO_[0], HIGH);
    digitalWrite(PO_[1], HIGH);
    digitalWrite(PO_[2], HIGH);
    digitalWrite(PO_[3], HIGH);
    digitalWrite(PO_[x], LOW);

    for (int y = 0; y < 16; y++)
    {
      bI[x][y] = digitalRead(PI_[y]);
      if (bI[x][y] == 0)
      {
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

void CheckMissed()
{
  bMissedDart = false;

  // Read both piezos
  for (int i = 0; i < 2; i++)
  {
    iValue[i] = analogRead(piezoPin[i]);
    iValue[i] = analogRead(piezoPin[i]);

    if (iValue[i] >= iPiezoThreshold)
    {
      bMissedDart = true;
    }
  }

  if (!bHitDetected && bMissedDart)
  {
    Serial.println("m");
    delay(300);
  }
}

void BlinkExtraFast(int times)
{
  for (int i = 0; i < times; i++)
  {
    digitalWrite(buttonLedPin, HIGH);
    delay(100);
    digitalWrite(buttonLedPin, LOW);
    delay(100);
  }
}

void Blink(int times)
{
  for (int i = 0; i < times; i++)
  {
    digitalWrite(buttonLedPin, HIGH);
    delay(250);
    digitalWrite(buttonLedPin, LOW);
    delay(250);
  }
}

void BlinkSlow(int times)
{
  for (int i = 0; i < times; i++)
  {
    digitalWrite(buttonLedPin, HIGH);
    delay(750);
    digitalWrite(buttonLedPin, LOW);
    delay(750);
  }
}

void BlinkExtraSlow(int times)
{
  for (int i = 0; i < times; i++)
  {
    digitalWrite(buttonLedPin, HIGH);
    delay(1000);
    digitalWrite(buttonLedPin, LOW);
    delay(1000);
  }
}

void SetUltrasonicThreshold()
{
  int iTempDistance = ReadUltrasonicDistance();
  iUltrasonicThreshold = iTempDistance + 3;
  bUltrasonicThresholdMeasured = true;
}

void ProcessUltrasonic()
{
  int iTempDistance = ReadUltrasonicDistance();
  if (iDistance > iUltrasonicThreshold)
  {
    Serial.println("u");
  }
}

void ProcessSerial()
{
  // Write different states with 0-5
  if (sInputString.indexOf("0") != -1)
  {
    // 0: INIT
    iState = 0;
    // Serial.println("DEBUG from ARD state set to 0");
  }
  if (sInputString.indexOf("1") != -1)
  {
    // 1: THROW
    iState = 1;
    ButtonOff();
    // Serial.println("DEBUG from ARD state set to 1");
  }
  if (sInputString.indexOf("2") != -1)
  {
    // 2: NEXTPLAYER
    iState = 2;
    // Serial.println("DEBUG from ARD state set to 2");
  }
  if (sInputString.indexOf("3") != -1)
  {
    // 3: MOTION DETECTED
    bMotionDetected = true;
    // Serial.println("DEBUG from ARD motion detected true");
  }
  if (sInputString.indexOf("4") != -1)
  {
    // 4: RESET ULTRASONIC
    bUltrasonicThresholdMeasured = false;
    bMotionDetected = false;
    ButtonOff();
    // Serial.println("DEBUG from ARD reset us");
  }
  if (sInputString.indexOf("5") != -1)
  {
    // 5: WON
    iState = 5;
    // Serial.println("DEBUG from ARD state set to 5");
  }
  // 6 - Button on
  if (sInputString.indexOf("6") != -1)
  {
    ButtonOn();
  }
  // 7 - Button off
  else if (sInputString.indexOf("7") != -1)
  {
    ButtonOff();
  }
  else if (sInputString.indexOf("9") != -1)
  {
    Blink(7);
    iState = 1;
  } else if (sInputString.startsWith("P"))
  {
    String value = sInputString.substring(1);
    iPiezoThreshold = value.toInt();
    // Serial.println("P: " + iPiezoThreshold);
  }

  sInputString = "";
  bStringComplete = false;
}

/* Setup loop */
void setup()
{
  // Pins Ultrasonic
  pinMode(trigPin, OUTPUT);
  pinMode(echoPin, INPUT);
  // Pin Button
  pinMode(buttonPin, INPUT_PULLUP);
  // Pin Button LED
  pinMode(buttonLedPin, OUTPUT);
  digitalWrite(buttonLedPin, LOW);
  // Matrix
  for (int i = 0; i < 4; i++)
  {
    pinMode(PO_[i], OUTPUT);
  }

  for (int i = 0; i < 16; i++)
  {
    pinMode(PI_[i], INPUT_PULLUP);
  }

  // Button blink 5 times
  Blink(5);
  // Start serial
  Serial.begin(9600);
}

/* Main loop */
void loop()
{
  // First read serial
  if (bStringComplete)
  {
    ProcessSerial();
  }

  if (iState == 0)
  {
    BlinkSlow(1);
  }
  else if (iState == 1)
  {
    EvalThrow();
    CheckMissed();
    CheckButton();
  }
  else if (iState == 2)
  {
    if (bUltrasonicThresholdMeasured)
    {
      CheckButton();
      // If there is no motion process ultrasonic
      if (!bMotionDetected)
      {
        ProcessUltrasonic();
      }
      // else blink to indicate dart retrieval
      else
      {
        Blink(1);
      }
    }
    else
    {
      ButtonOn();
      SetUltrasonicThreshold();
    }
  }
  else if (iState == 5)
  {
    CheckButton();
    BlinkExtraSlow(1);
  }
}

/* Serial Events */
void serialEvent()
{
  while (Serial.available())
  {
    char inChar = (char)Serial.read();
    sInputString += inChar;
    if (inChar == '\n')
    {
      bStringComplete = true;
    }
  }
}
