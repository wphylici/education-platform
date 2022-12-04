const jwt = require("jsonwebtoken");

// const { NODE_ENV, JWT_SECRET } = process.env;
// const UnauthorizedError = require("../errors/unauthorized-err");
// const { MESSAGE_ERROR_401 } = require("../utils/constants");

module.exports = (req, res, next) => {
  const { authorization } = req.headers;

  if (!authorization || !authorization.startsWith("Bearer ")) {
    return res.status(401).send({ message: "Необходима авторизация" });
    // throw new UnauthorizedError(MESSAGE_ERROR_401);
  }

  const token = authorization.replace("Bearer ", "");

  let payload;

  try {
    payload = jwt.verify(
      token,
      "some-secret-key"
      // NODE_ENV === "production" ? JWT_SECRET : "some-secret-key"
    );

    console.log(payload);
  } catch (err) {
    return res.status(401).send({ message: "Необходима авторизация" });
    // throw new UnauthorizedError(MESSAGE_ERROR_401);
  }

  req.user = payload;

  next();
};
