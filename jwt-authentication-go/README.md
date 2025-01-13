# Creating a Simple ToDo Application

For the purpose of demonstrating JWT authentication, we will create a simple ToDo list application with [Role-Based Access Control (RBAC)](https://www.upguard.com/blog/rbac). We will introduce two distinct roles: senior and employee. These roles will serve as the basis for managing user access and permissions within our application.

  * Senior: Senior users have elevated privileges and can perform actions such as adding ToDo items to the list.

  * Employee: Employees have more restricted access and cannot add new ToDo items.

Roles are typically assigned during user authentication and embedded in the JWT claims.

Now, let's dive into the fun part – creating our basic ToDo application using the powerful [Gin framework.](https://gin-gonic.com/) This section will walk you through the steps, breaking down the code into manageable snippets.