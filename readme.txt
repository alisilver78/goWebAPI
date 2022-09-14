User struct:
{
	Email    string 
	Password string 
	IsAdmin  bool
}
// example 
// {
// "admin@example.com"
// "12345abc"
// true
// }
product struct:
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
// example 
// {
// "iphone13"
// 599
// "USD"
// 30
// 0.05
// "apple"
// ["cable","manuals"]
// }

requared access:
DELETE: isAdmin = true, jwd token //header: "x-auth-token"
GET: none
POST: jwd token
PUT: jwd token

Users endpoint /users creates user
Auth endpoint /auth authoriz an existing user

///////////////////////tasks////////////////////////////
create an admin user and add adminMiddleware to user endpoint