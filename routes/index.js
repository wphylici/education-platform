const router = require("express").Router();
const userRoutes = require("./user");
const courseRoutes = require("./course")

router.use(userRoutes, courseRoutes);

module.exports = router;
