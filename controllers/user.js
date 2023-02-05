const bcrypt = require("bcryptjs");
const jwt = require("jsonwebtoken"); // для создания токена
const User = require("../models/user");
const {
  MESSAGE_ERROR_404,
  MESSAGE_ERROR_400,
  MESSAGE_ERROR_409,
} = require("../utils/constants");

const NotFoundError = require("../errors/not-found-err");
const BadRequestError = require("../errors/bad-request-err");
const ConflictError = require("../errors/conflict-err");

// регистрация
createUser = (req, res, next) => {
  const { name, surname, fathername, email, password, status, groups } =
    req.body;

  bcrypt.hash(password, 10).then((hash) => {
    User.create({
      name,
      surname,
      fathername,
      email,
      password: hash,
      status,
      groups,
    })
      .then((user) =>
        res.send({ _id: user._id, name: user.name, email: user.email })
      )
      .catch((err) => {
        if (err.message === "CastError" || err.message === "ValidationError") {
          throw new BadRequestError(MESSAGE_ERROR_400);
        } else if (err.code === 11000) {
          throw new ConflictError(MESSAGE_ERROR_409);
        }
        console.log(err);
        throw err;
      })
      .catch(next);
  });
};

//вход
const login = (req, res, next) => {
  const { email, password } = req.body;

  return User.findUserByCredentials(email, password)
    .then((user) => {
      const token = jwt.sign(
        { _id: user._id },
        "some-secret-key",
        // NODE_ENV === "production" ? JWT_SECRET : "some-secret-key",
        {
          expiresIn: "7d",
        }
      );
      res.send({ token });
    })
    .catch((err) => {
      res.status(401).send({ message: err.message });
      throw new BadRequestError(MESSAGE_ERROR_400);
    })
    .catch(next);
};

// возвращение текущего пользователя
const returnCurrentUser = (req, res, next) => {
  User.findById(req.user._id)
    .then((user) => {
      if (!user) {
        console.log("пользователь не найден");
        return next(new NotFoundError(MESSAGE_ERROR_404));
      }
      return res.send({
        name: user.name,
        surname: user.surname,
        fathername: user.fathername,
        email: user.email,
        status: user.status,
        groups: user.groups,
      });
    })
    .catch(next);
};

module.exports = { createUser, login, returnCurrentUser };
