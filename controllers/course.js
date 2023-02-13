const Course = require('../models/course')
// const BadRequestError = require("../errors/bad-request-err");
// const {MESSAGE_ERROR_400, MESSAGE_ERROR_409} = require("../utils/constants");
// const ConflictError = require("../errors/conflict-err");

createCourse = (req, res, next) => {
  const { name, createdBy, chapters } = req.body;

  Course.create({name, createdBy, chapters})
    .then((course) => res.send(course))
    .catch(err => res.status(500).send({ message: 'Произошла ошибка' }));


};

getCourse = (req, res, next) => {
  const { createdBy } = req.body;
    Course.find({ createdBy })
      .then(course => res.send(course))
      .catch(err => res.status(500).send({ message: 'Произошла ошибка' }));
}

module.exports = { createCourse, getCourse };