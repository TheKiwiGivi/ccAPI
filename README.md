Heroku: https://polar-hollows-99028.herokuapp.com/country

This program has six functions.
1. It can find out if two countries are in the same region.
2. It can find out if two countries are connected directly by border (does not apply to sea border meaning that countries like Japan does not have any bordering countries)
3. It can find out which of the two countries has the highest population
4. If population is requested, it will add all new countries to a mongodb.
5. Entries to the database can be removed by the /population/remove/<countryname>
6. There is a ranking system that shows the country with the highest population in the database, accessed by /population/ranks/

For region add /region
For border add /border
For population add /population
	To remove a country from the database add /<countryname>
It utilizes a free API from https://restcountries.eu which was shown in class. 

Further ideas would be to get the currency name from the API and use another API to also be able to compare currencies in a country to determin the strongest one. 

The API does not require to type in the whole name of the country, of you omit some parts it will choose the country closest to what was typed. 

Only GET method is required.

Please always end a URL with a '/' to make sure everything will work properly. 
 * The only exception to this is when removing a country, here, the URL should NOT end in "/"*

Some names can only be accessed by using alternative spellings, which can be found on the website when calling a specific country like https://restcountries.eu/rest/v2/name/norway

It should be noted that when removing a country, it has to be types in exactly as it is called in the population output.
EXAMPLES: 
https://polar-hollows-99028.herokuapp.com/country/border/peru/brazil/
https://polar-hollows-99028.herokuapp.com/country/region/spain/morocco/
https://polar-hollows-99028.herokuapp.com/country/population/norway/sweden/
https://polar-hollows-99028.herokuapp.com/country/population/remove/Sweden
https://polar-hollows-99028.herokuapp.com/country/population/ranks/

