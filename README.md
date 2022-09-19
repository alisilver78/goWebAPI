## Document fields

User struct: //located in users.go
{
	Email    string 
	Password string 
	IsAdmin  bool
}
###### example 
	{
	"email": "admin@example.com",
	"password": "12345abc",
	"isadmin": true
	}
###### product struct: (located in products.go)
	{
	ID          primitive.ObjectID 
	Name        string             
	Price       int                
	Currency    string             
	Quantity    int                
	Discount    int                
	Vendor      string             
	Accessories []string           
	}
###### example 
	{
	"product_name": "iphone13",
	"price": 599,
	"currency": "USD",
	"quantity": 30,
	"discount": 0.05,
	"vendor": "apple",
	"accessories": ["cable","manuals"]
	}

## requared access for products endpoind methods
DELETE: isAdmin = true, jwd token (header: "x-auth-token" for jwd token)
POST: jwd token (header: "x-auth-token" for jwd token)
PUT: jwd token (header: "x-auth-token" for jwd token)
GET: none


## Endpoints

Users endpoint /users creates user
Auth endpoint /auth authoriz an existing user


## Files

main.go: endpoints, middlewares
products.go: products endpoint's handlers and struct
users.go: users endpoint's handlers and struct
validators.go: validator structs and methods for both users and products
config.go: config properties and variable
interface.go: CollectionAPI interface

------------------------------ **remaining tasks** ---------------------------------
()create an admin user and add adminMiddleware to user endpoint
