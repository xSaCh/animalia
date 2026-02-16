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

## Behavior Trees
If thirsty find water 
if hungry find food 
if tired find rest Area
    rest for a while
Find Roam Pos
    Roam to that pos

Root (Selector)
├─ Sequence: (Thirst) -> Find Water
│   ├─ Condition: is_thirsty?
│   └─ Action: find_water
└─ Sequence: (Hunger) -> Find Food
│   ├─ Condition: is_hungry?
│   └─ Action: find_food
└─ Sequence: (Tiredness) -> Rest
│   ├─ Condition: is_tired?
│   ├─ Condition: find_rest_area?
│   └─ Action: rest
└─ Sequence: (Roam) -> Roam
    ├─ Condition: find_roam_pos?
    └─ Action: roam_to_pos