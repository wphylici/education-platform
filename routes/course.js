const courseRoutes = require("express").Router();
const { createCourse, getCourse } = require("../controllers/course");

courseRoutes.post("/course", createCourse);
courseRoutes.get("/course", getCourse)
module.exports = courseRoutes;
