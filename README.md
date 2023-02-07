# Recipe Stats Calculator

Application reads JSON files*, calculates data and returns JSON file with final results. Below you will find detailed description.

---
*Example of input JSON file

```
[
 ...
 {
    "postcode": "10224",
    "recipe": "Creamy Dill Chicken",
    "delivery": "Wednesday 1AM - 7PM"
 },
 ...
]
```
---
## Application setup
1. Clone git to your machine
2. Create `.env` file as in example and provide info (input file, file path, output file path)
3. Run build app or run `main.go`
4. Provide info for calculations
5. Done. Find output file in output directory
---
## App functional description
- Has CLI to provide user data for calculations
- Counts the number of unique recipe names
- Counts the number of occurences for each unique recipe name (alphabetically ordered by recipe name)
- Finds the postcode with most delivered recipes
- Counts the number of deliveries to given postcode that lie within the given delivery time
- Lists the recipe names (alphabetically ordered) that contain in their name one of the given words
- Expected output is rendered to `stdout` folder (can be cofigured)