const userRoutes = require("express").Router();
const { createUser, login, returnCurrentUser } = require("../controllers/user");
const auth = require("../middlewares/auth");

userRoutes.post("/signup", createUser);
userRoutes.post("/signin", login);

userRoutes.use(auth);

userRoutes.get("/users/me", returnCurrentUser);

module.exports = userRoutes;
