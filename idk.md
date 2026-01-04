# Entity

## State Machine

    Priority: Thirst > Hunger > Tiredness > Roam

    Next State:
        IF Thirst >= 80: FindWater

        IF Hunger >= 70: FindFood

        IF Tiredness >= 85: Rest

        // Exit conditions
        IF currentState == FindWater AND Thirst <= 20:  Roam

        IF currentState == FindFood AND Hunger <= 25:  Roam

        IF currentState == Rest AND Tiredness <= 30:  Roam

    Stats:
        Roam:
            Hunger    += 1
            Thirst    += 2
            Tiredness += 1
            BREAK

        FindFood:
            Hunger    -= 3        // eating
            Thirst    += 1
            Tiredness += 1
            BREAK

        FindWater:
            Thirst    -= 4        // drinking
            Hunger    += 1
            Tiredness += 1
            BREAK

        Rest:
            Tiredness -= 4
            Hunger    += 1
            Thirst    += 1
            BREAK
