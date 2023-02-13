const mongoose = require("mongoose");
// const validator = require("validator");

const courseSchema = new mongoose.Schema(
  {
    createdBy: {
      type: mongoose.Schema.Types.ObjectId,
      ref: 'user',
      required: true,
    },
    chapters: {
      type: Array,
      required: true,
    },
    name: {
      type: String,
      required: true,
    }
  },
  { versionKey: false }
);

module.exports = mongoose.model("Course", courseSchema);
