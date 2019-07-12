/*eslint-disable block-scoped-var, id-length, no-control-regex, no-magic-numbers, no-prototype-builtins, no-redeclare, no-shadow, no-var, sort-vars*/
"use strict";

var $protobuf = require("protobufjs/minimal");

// Common aliases
var $Reader = $protobuf.Reader, $Writer = $protobuf.Writer, $util = $protobuf.util;

// Exported root namespace
var $root = $protobuf.roots["default"] || ($protobuf.roots["default"] = {});

$root.message = (function() {

    /**
     * Namespace message.
     * @exports message
     * @namespace
     */
    var message = {};

    message.Point3F = (function() {

        /**
         * Properties of a Point3F.
         * @memberof message
         * @interface IPoint3F
         * @property {number|null} [X] Point3F X
         * @property {number|null} [Y] Point3F Y
         * @property {number|null} [Z] Point3F Z
         */

        /**
         * Constructs a new Point3F.
         * @memberof message
         * @classdesc Represents a Point3F.
         * @implements IPoint3F
         * @constructor
         * @param {message.IPoint3F=} [properties] Properties to set
         */
        function Point3F(properties) {
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * Point3F X.
         * @member {number} X
         * @memberof message.Point3F
         * @instance
         */
        Point3F.prototype.X = 0;

        /**
         * Point3F Y.
         * @member {number} Y
         * @memberof message.Point3F
         * @instance
         */
        Point3F.prototype.Y = 0;

        /**
         * Point3F Z.
         * @member {number} Z
         * @memberof message.Point3F
         * @instance
         */
        Point3F.prototype.Z = 0;

        /**
         * Creates a new Point3F instance using the specified properties.
         * @function create
         * @memberof message.Point3F
         * @static
         * @param {message.IPoint3F=} [properties] Properties to set
         * @returns {message.Point3F} Point3F instance
         */
        Point3F.create = function create(properties) {
            return new Point3F(properties);
        };

        /**
         * Encodes the specified Point3F message. Does not implicitly {@link message.Point3F.verify|verify} messages.
         * @function encode
         * @memberof message.Point3F
         * @static
         * @param {message.IPoint3F} message Point3F message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        Point3F.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.X != null && message.hasOwnProperty("X"))
                writer.uint32(/* id 1, wireType 5 =*/13).float(message.X);
            if (message.Y != null && message.hasOwnProperty("Y"))
                writer.uint32(/* id 2, wireType 5 =*/21).float(message.Y);
            if (message.Z != null && message.hasOwnProperty("Z"))
                writer.uint32(/* id 3, wireType 5 =*/29).float(message.Z);
            return writer;
        };

        /**
         * Encodes the specified Point3F message, length delimited. Does not implicitly {@link message.Point3F.verify|verify} messages.
         * @function encodeDelimited
         * @memberof message.Point3F
         * @static
         * @param {message.IPoint3F} message Point3F message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        Point3F.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a Point3F message from the specified reader or buffer.
         * @function decode
         * @memberof message.Point3F
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {message.Point3F} Point3F
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        Point3F.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.Point3F();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.X = reader.float();
                    break;
                case 2:
                    message.Y = reader.float();
                    break;
                case 3:
                    message.Z = reader.float();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a Point3F message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof message.Point3F
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {message.Point3F} Point3F
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        Point3F.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a Point3F message.
         * @function verify
         * @memberof message.Point3F
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        Point3F.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.X != null && message.hasOwnProperty("X"))
                if (typeof message.X !== "number")
                    return "X: number expected";
            if (message.Y != null && message.hasOwnProperty("Y"))
                if (typeof message.Y !== "number")
                    return "Y: number expected";
            if (message.Z != null && message.hasOwnProperty("Z"))
                if (typeof message.Z !== "number")
                    return "Z: number expected";
            return null;
        };

        /**
         * Creates a Point3F message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof message.Point3F
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {message.Point3F} Point3F
         */
        Point3F.fromObject = function fromObject(object) {
            if (object instanceof $root.message.Point3F)
                return object;
            var message = new $root.message.Point3F();
            if (object.X != null)
                message.X = Number(object.X);
            if (object.Y != null)
                message.Y = Number(object.Y);
            if (object.Z != null)
                message.Z = Number(object.Z);
            return message;
        };

        /**
         * Creates a plain object from a Point3F message. Also converts values to other types if specified.
         * @function toObject
         * @memberof message.Point3F
         * @static
         * @param {message.Point3F} message Point3F
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        Point3F.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.defaults) {
                object.X = 0;
                object.Y = 0;
                object.Z = 0;
            }
            if (message.X != null && message.hasOwnProperty("X"))
                object.X = options.json && !isFinite(message.X) ? String(message.X) : message.X;
            if (message.Y != null && message.hasOwnProperty("Y"))
                object.Y = options.json && !isFinite(message.Y) ? String(message.Y) : message.Y;
            if (message.Z != null && message.hasOwnProperty("Z"))
                object.Z = options.json && !isFinite(message.Z) ? String(message.Z) : message.Z;
            return object;
        };

        /**
         * Converts this Point3F to JSON.
         * @function toJSON
         * @memberof message.Point3F
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        Point3F.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return Point3F;
    })();

    message.C_W_Move = (function() {

        /**
         * Properties of a C_W_Move.
         * @memberof message
         * @interface IC_W_Move
         * @property {message.IIpacket|null} [PacketHead] C_W_Move PacketHead
         * @property {message.C_W_Move.IMove|null} [move] C_W_Move move
         */

        /**
         * Constructs a new C_W_Move.
         * @memberof message
         * @classdesc Represents a C_W_Move.
         * @implements IC_W_Move
         * @constructor
         * @param {message.IC_W_Move=} [properties] Properties to set
         */
        function C_W_Move(properties) {
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * C_W_Move PacketHead.
         * @member {message.IIpacket|null|undefined} PacketHead
         * @memberof message.C_W_Move
         * @instance
         */
        C_W_Move.prototype.PacketHead = null;

        /**
         * C_W_Move move.
         * @member {message.C_W_Move.IMove|null|undefined} move
         * @memberof message.C_W_Move
         * @instance
         */
        C_W_Move.prototype.move = null;

        /**
         * Creates a new C_W_Move instance using the specified properties.
         * @function create
         * @memberof message.C_W_Move
         * @static
         * @param {message.IC_W_Move=} [properties] Properties to set
         * @returns {message.C_W_Move} C_W_Move instance
         */
        C_W_Move.create = function create(properties) {
            return new C_W_Move(properties);
        };

        /**
         * Encodes the specified C_W_Move message. Does not implicitly {@link message.C_W_Move.verify|verify} messages.
         * @function encode
         * @memberof message.C_W_Move
         * @static
         * @param {message.IC_W_Move} message C_W_Move message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        C_W_Move.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                $root.message.Ipacket.encode(message.PacketHead, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
            if (message.move != null && message.hasOwnProperty("move"))
                $root.message.C_W_Move.Move.encode(message.move, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
            return writer;
        };

        /**
         * Encodes the specified C_W_Move message, length delimited. Does not implicitly {@link message.C_W_Move.verify|verify} messages.
         * @function encodeDelimited
         * @memberof message.C_W_Move
         * @static
         * @param {message.IC_W_Move} message C_W_Move message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        C_W_Move.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a C_W_Move message from the specified reader or buffer.
         * @function decode
         * @memberof message.C_W_Move
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {message.C_W_Move} C_W_Move
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        C_W_Move.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.C_W_Move();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.PacketHead = $root.message.Ipacket.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.move = $root.message.C_W_Move.Move.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a C_W_Move message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof message.C_W_Move
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {message.C_W_Move} C_W_Move
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        C_W_Move.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a C_W_Move message.
         * @function verify
         * @memberof message.C_W_Move
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        C_W_Move.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead")) {
                var error = $root.message.Ipacket.verify(message.PacketHead);
                if (error)
                    return "PacketHead." + error;
            }
            if (message.move != null && message.hasOwnProperty("move")) {
                var error = $root.message.C_W_Move.Move.verify(message.move);
                if (error)
                    return "move." + error;
            }
            return null;
        };

        /**
         * Creates a C_W_Move message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof message.C_W_Move
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {message.C_W_Move} C_W_Move
         */
        C_W_Move.fromObject = function fromObject(object) {
            if (object instanceof $root.message.C_W_Move)
                return object;
            var message = new $root.message.C_W_Move();
            if (object.PacketHead != null) {
                if (typeof object.PacketHead !== "object")
                    throw TypeError(".message.C_W_Move.PacketHead: object expected");
                message.PacketHead = $root.message.Ipacket.fromObject(object.PacketHead);
            }
            if (object.move != null) {
                if (typeof object.move !== "object")
                    throw TypeError(".message.C_W_Move.move: object expected");
                message.move = $root.message.C_W_Move.Move.fromObject(object.move);
            }
            return message;
        };

        /**
         * Creates a plain object from a C_W_Move message. Also converts values to other types if specified.
         * @function toObject
         * @memberof message.C_W_Move
         * @static
         * @param {message.C_W_Move} message C_W_Move
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        C_W_Move.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.defaults) {
                object.PacketHead = null;
                object.move = null;
            }
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                object.PacketHead = $root.message.Ipacket.toObject(message.PacketHead, options);
            if (message.move != null && message.hasOwnProperty("move"))
                object.move = $root.message.C_W_Move.Move.toObject(message.move, options);
            return object;
        };

        /**
         * Converts this C_W_Move to JSON.
         * @function toJSON
         * @memberof message.C_W_Move
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        C_W_Move.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        C_W_Move.Move = (function() {

            /**
             * Properties of a Move.
             * @memberof message.C_W_Move
             * @interface IMove
             * @property {number|null} [Mode] Move Mode
             * @property {message.C_W_Move.Move.INormal|null} [normal] Move normal
             * @property {message.C_W_Move.Move.IPath|null} [path] Move path
             * @property {message.C_W_Move.Move.IBlink|null} [link] Move link
             * @property {message.C_W_Move.Move.IJump|null} [jump] Move jump
             * @property {message.C_W_Move.Move.ILine|null} [line] Move line
             */

            /**
             * Constructs a new Move.
             * @memberof message.C_W_Move
             * @classdesc Represents a Move.
             * @implements IMove
             * @constructor
             * @param {message.C_W_Move.IMove=} [properties] Properties to set
             */
            function Move(properties) {
                if (properties)
                    for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * Move Mode.
             * @member {number} Mode
             * @memberof message.C_W_Move.Move
             * @instance
             */
            Move.prototype.Mode = 0;

            /**
             * Move normal.
             * @member {message.C_W_Move.Move.INormal|null|undefined} normal
             * @memberof message.C_W_Move.Move
             * @instance
             */
            Move.prototype.normal = null;

            /**
             * Move path.
             * @member {message.C_W_Move.Move.IPath|null|undefined} path
             * @memberof message.C_W_Move.Move
             * @instance
             */
            Move.prototype.path = null;

            /**
             * Move link.
             * @member {message.C_W_Move.Move.IBlink|null|undefined} link
             * @memberof message.C_W_Move.Move
             * @instance
             */
            Move.prototype.link = null;

            /**
             * Move jump.
             * @member {message.C_W_Move.Move.IJump|null|undefined} jump
             * @memberof message.C_W_Move.Move
             * @instance
             */
            Move.prototype.jump = null;

            /**
             * Move line.
             * @member {message.C_W_Move.Move.ILine|null|undefined} line
             * @memberof message.C_W_Move.Move
             * @instance
             */
            Move.prototype.line = null;

            /**
             * Creates a new Move instance using the specified properties.
             * @function create
             * @memberof message.C_W_Move.Move
             * @static
             * @param {message.C_W_Move.IMove=} [properties] Properties to set
             * @returns {message.C_W_Move.Move} Move instance
             */
            Move.create = function create(properties) {
                return new Move(properties);
            };

            /**
             * Encodes the specified Move message. Does not implicitly {@link message.C_W_Move.Move.verify|verify} messages.
             * @function encode
             * @memberof message.C_W_Move.Move
             * @static
             * @param {message.C_W_Move.IMove} message Move message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            Move.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.Mode != null && message.hasOwnProperty("Mode"))
                    writer.uint32(/* id 1, wireType 0 =*/8).int32(message.Mode);
                if (message.normal != null && message.hasOwnProperty("normal"))
                    $root.message.C_W_Move.Move.Normal.encode(message.normal, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
                if (message.path != null && message.hasOwnProperty("path"))
                    $root.message.C_W_Move.Move.Path.encode(message.path, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
                if (message.link != null && message.hasOwnProperty("link"))
                    $root.message.C_W_Move.Move.Blink.encode(message.link, writer.uint32(/* id 4, wireType 2 =*/34).fork()).ldelim();
                if (message.jump != null && message.hasOwnProperty("jump"))
                    $root.message.C_W_Move.Move.Jump.encode(message.jump, writer.uint32(/* id 5, wireType 2 =*/42).fork()).ldelim();
                if (message.line != null && message.hasOwnProperty("line"))
                    $root.message.C_W_Move.Move.Line.encode(message.line, writer.uint32(/* id 6, wireType 2 =*/50).fork()).ldelim();
                return writer;
            };

            /**
             * Encodes the specified Move message, length delimited. Does not implicitly {@link message.C_W_Move.Move.verify|verify} messages.
             * @function encodeDelimited
             * @memberof message.C_W_Move.Move
             * @static
             * @param {message.C_W_Move.IMove} message Move message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            Move.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };

            /**
             * Decodes a Move message from the specified reader or buffer.
             * @function decode
             * @memberof message.C_W_Move.Move
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {message.C_W_Move.Move} Move
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            Move.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.C_W_Move.Move();
                while (reader.pos < end) {
                    var tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1:
                        message.Mode = reader.int32();
                        break;
                    case 2:
                        message.normal = $root.message.C_W_Move.Move.Normal.decode(reader, reader.uint32());
                        break;
                    case 3:
                        message.path = $root.message.C_W_Move.Move.Path.decode(reader, reader.uint32());
                        break;
                    case 4:
                        message.link = $root.message.C_W_Move.Move.Blink.decode(reader, reader.uint32());
                        break;
                    case 5:
                        message.jump = $root.message.C_W_Move.Move.Jump.decode(reader, reader.uint32());
                        break;
                    case 6:
                        message.line = $root.message.C_W_Move.Move.Line.decode(reader, reader.uint32());
                        break;
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };

            /**
             * Decodes a Move message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof message.C_W_Move.Move
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {message.C_W_Move.Move} Move
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            Move.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies a Move message.
             * @function verify
             * @memberof message.C_W_Move.Move
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            Move.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.Mode != null && message.hasOwnProperty("Mode"))
                    if (!$util.isInteger(message.Mode))
                        return "Mode: integer expected";
                if (message.normal != null && message.hasOwnProperty("normal")) {
                    var error = $root.message.C_W_Move.Move.Normal.verify(message.normal);
                    if (error)
                        return "normal." + error;
                }
                if (message.path != null && message.hasOwnProperty("path")) {
                    var error = $root.message.C_W_Move.Move.Path.verify(message.path);
                    if (error)
                        return "path." + error;
                }
                if (message.link != null && message.hasOwnProperty("link")) {
                    var error = $root.message.C_W_Move.Move.Blink.verify(message.link);
                    if (error)
                        return "link." + error;
                }
                if (message.jump != null && message.hasOwnProperty("jump")) {
                    var error = $root.message.C_W_Move.Move.Jump.verify(message.jump);
                    if (error)
                        return "jump." + error;
                }
                if (message.line != null && message.hasOwnProperty("line")) {
                    var error = $root.message.C_W_Move.Move.Line.verify(message.line);
                    if (error)
                        return "line." + error;
                }
                return null;
            };

            /**
             * Creates a Move message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof message.C_W_Move.Move
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {message.C_W_Move.Move} Move
             */
            Move.fromObject = function fromObject(object) {
                if (object instanceof $root.message.C_W_Move.Move)
                    return object;
                var message = new $root.message.C_W_Move.Move();
                if (object.Mode != null)
                    message.Mode = object.Mode | 0;
                if (object.normal != null) {
                    if (typeof object.normal !== "object")
                        throw TypeError(".message.C_W_Move.Move.normal: object expected");
                    message.normal = $root.message.C_W_Move.Move.Normal.fromObject(object.normal);
                }
                if (object.path != null) {
                    if (typeof object.path !== "object")
                        throw TypeError(".message.C_W_Move.Move.path: object expected");
                    message.path = $root.message.C_W_Move.Move.Path.fromObject(object.path);
                }
                if (object.link != null) {
                    if (typeof object.link !== "object")
                        throw TypeError(".message.C_W_Move.Move.link: object expected");
                    message.link = $root.message.C_W_Move.Move.Blink.fromObject(object.link);
                }
                if (object.jump != null) {
                    if (typeof object.jump !== "object")
                        throw TypeError(".message.C_W_Move.Move.jump: object expected");
                    message.jump = $root.message.C_W_Move.Move.Jump.fromObject(object.jump);
                }
                if (object.line != null) {
                    if (typeof object.line !== "object")
                        throw TypeError(".message.C_W_Move.Move.line: object expected");
                    message.line = $root.message.C_W_Move.Move.Line.fromObject(object.line);
                }
                return message;
            };

            /**
             * Creates a plain object from a Move message. Also converts values to other types if specified.
             * @function toObject
             * @memberof message.C_W_Move.Move
             * @static
             * @param {message.C_W_Move.Move} message Move
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            Move.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                var object = {};
                if (options.defaults) {
                    object.Mode = 0;
                    object.normal = null;
                    object.path = null;
                    object.link = null;
                    object.jump = null;
                    object.line = null;
                }
                if (message.Mode != null && message.hasOwnProperty("Mode"))
                    object.Mode = message.Mode;
                if (message.normal != null && message.hasOwnProperty("normal"))
                    object.normal = $root.message.C_W_Move.Move.Normal.toObject(message.normal, options);
                if (message.path != null && message.hasOwnProperty("path"))
                    object.path = $root.message.C_W_Move.Move.Path.toObject(message.path, options);
                if (message.link != null && message.hasOwnProperty("link"))
                    object.link = $root.message.C_W_Move.Move.Blink.toObject(message.link, options);
                if (message.jump != null && message.hasOwnProperty("jump"))
                    object.jump = $root.message.C_W_Move.Move.Jump.toObject(message.jump, options);
                if (message.line != null && message.hasOwnProperty("line"))
                    object.line = $root.message.C_W_Move.Move.Line.toObject(message.line, options);
                return object;
            };

            /**
             * Converts this Move to JSON.
             * @function toJSON
             * @memberof message.C_W_Move.Move
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            Move.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            Move.Normal = (function() {

                /**
                 * Properties of a Normal.
                 * @memberof message.C_W_Move.Move
                 * @interface INormal
                 * @property {message.IPoint3F|null} [Pos] Normal Pos
                 * @property {number|null} [Yaw] Normal Yaw
                 * @property {number|null} [Duration] Normal Duration
                 */

                /**
                 * Constructs a new Normal.
                 * @memberof message.C_W_Move.Move
                 * @classdesc Represents a Normal.
                 * @implements INormal
                 * @constructor
                 * @param {message.C_W_Move.Move.INormal=} [properties] Properties to set
                 */
                function Normal(properties) {
                    if (properties)
                        for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                            if (properties[keys[i]] != null)
                                this[keys[i]] = properties[keys[i]];
                }

                /**
                 * Normal Pos.
                 * @member {message.IPoint3F|null|undefined} Pos
                 * @memberof message.C_W_Move.Move.Normal
                 * @instance
                 */
                Normal.prototype.Pos = null;

                /**
                 * Normal Yaw.
                 * @member {number} Yaw
                 * @memberof message.C_W_Move.Move.Normal
                 * @instance
                 */
                Normal.prototype.Yaw = 0;

                /**
                 * Normal Duration.
                 * @member {number} Duration
                 * @memberof message.C_W_Move.Move.Normal
                 * @instance
                 */
                Normal.prototype.Duration = 0;

                /**
                 * Creates a new Normal instance using the specified properties.
                 * @function create
                 * @memberof message.C_W_Move.Move.Normal
                 * @static
                 * @param {message.C_W_Move.Move.INormal=} [properties] Properties to set
                 * @returns {message.C_W_Move.Move.Normal} Normal instance
                 */
                Normal.create = function create(properties) {
                    return new Normal(properties);
                };

                /**
                 * Encodes the specified Normal message. Does not implicitly {@link message.C_W_Move.Move.Normal.verify|verify} messages.
                 * @function encode
                 * @memberof message.C_W_Move.Move.Normal
                 * @static
                 * @param {message.C_W_Move.Move.INormal} message Normal message or plain object to encode
                 * @param {$protobuf.Writer} [writer] Writer to encode to
                 * @returns {$protobuf.Writer} Writer
                 */
                Normal.encode = function encode(message, writer) {
                    if (!writer)
                        writer = $Writer.create();
                    if (message.Pos != null && message.hasOwnProperty("Pos"))
                        $root.message.Point3F.encode(message.Pos, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
                    if (message.Yaw != null && message.hasOwnProperty("Yaw"))
                        writer.uint32(/* id 2, wireType 5 =*/21).float(message.Yaw);
                    if (message.Duration != null && message.hasOwnProperty("Duration"))
                        writer.uint32(/* id 3, wireType 5 =*/29).float(message.Duration);
                    return writer;
                };

                /**
                 * Encodes the specified Normal message, length delimited. Does not implicitly {@link message.C_W_Move.Move.Normal.verify|verify} messages.
                 * @function encodeDelimited
                 * @memberof message.C_W_Move.Move.Normal
                 * @static
                 * @param {message.C_W_Move.Move.INormal} message Normal message or plain object to encode
                 * @param {$protobuf.Writer} [writer] Writer to encode to
                 * @returns {$protobuf.Writer} Writer
                 */
                Normal.encodeDelimited = function encodeDelimited(message, writer) {
                    return this.encode(message, writer).ldelim();
                };

                /**
                 * Decodes a Normal message from the specified reader or buffer.
                 * @function decode
                 * @memberof message.C_W_Move.Move.Normal
                 * @static
                 * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
                 * @param {number} [length] Message length if known beforehand
                 * @returns {message.C_W_Move.Move.Normal} Normal
                 * @throws {Error} If the payload is not a reader or valid buffer
                 * @throws {$protobuf.util.ProtocolError} If required fields are missing
                 */
                Normal.decode = function decode(reader, length) {
                    if (!(reader instanceof $Reader))
                        reader = $Reader.create(reader);
                    var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.C_W_Move.Move.Normal();
                    while (reader.pos < end) {
                        var tag = reader.uint32();
                        switch (tag >>> 3) {
                        case 1:
                            message.Pos = $root.message.Point3F.decode(reader, reader.uint32());
                            break;
                        case 2:
                            message.Yaw = reader.float();
                            break;
                        case 3:
                            message.Duration = reader.float();
                            break;
                        default:
                            reader.skipType(tag & 7);
                            break;
                        }
                    }
                    return message;
                };

                /**
                 * Decodes a Normal message from the specified reader or buffer, length delimited.
                 * @function decodeDelimited
                 * @memberof message.C_W_Move.Move.Normal
                 * @static
                 * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
                 * @returns {message.C_W_Move.Move.Normal} Normal
                 * @throws {Error} If the payload is not a reader or valid buffer
                 * @throws {$protobuf.util.ProtocolError} If required fields are missing
                 */
                Normal.decodeDelimited = function decodeDelimited(reader) {
                    if (!(reader instanceof $Reader))
                        reader = new $Reader(reader);
                    return this.decode(reader, reader.uint32());
                };

                /**
                 * Verifies a Normal message.
                 * @function verify
                 * @memberof message.C_W_Move.Move.Normal
                 * @static
                 * @param {Object.<string,*>} message Plain object to verify
                 * @returns {string|null} `null` if valid, otherwise the reason why it is not
                 */
                Normal.verify = function verify(message) {
                    if (typeof message !== "object" || message === null)
                        return "object expected";
                    if (message.Pos != null && message.hasOwnProperty("Pos")) {
                        var error = $root.message.Point3F.verify(message.Pos);
                        if (error)
                            return "Pos." + error;
                    }
                    if (message.Yaw != null && message.hasOwnProperty("Yaw"))
                        if (typeof message.Yaw !== "number")
                            return "Yaw: number expected";
                    if (message.Duration != null && message.hasOwnProperty("Duration"))
                        if (typeof message.Duration !== "number")
                            return "Duration: number expected";
                    return null;
                };

                /**
                 * Creates a Normal message from a plain object. Also converts values to their respective internal types.
                 * @function fromObject
                 * @memberof message.C_W_Move.Move.Normal
                 * @static
                 * @param {Object.<string,*>} object Plain object
                 * @returns {message.C_W_Move.Move.Normal} Normal
                 */
                Normal.fromObject = function fromObject(object) {
                    if (object instanceof $root.message.C_W_Move.Move.Normal)
                        return object;
                    var message = new $root.message.C_W_Move.Move.Normal();
                    if (object.Pos != null) {
                        if (typeof object.Pos !== "object")
                            throw TypeError(".message.C_W_Move.Move.Normal.Pos: object expected");
                        message.Pos = $root.message.Point3F.fromObject(object.Pos);
                    }
                    if (object.Yaw != null)
                        message.Yaw = Number(object.Yaw);
                    if (object.Duration != null)
                        message.Duration = Number(object.Duration);
                    return message;
                };

                /**
                 * Creates a plain object from a Normal message. Also converts values to other types if specified.
                 * @function toObject
                 * @memberof message.C_W_Move.Move.Normal
                 * @static
                 * @param {message.C_W_Move.Move.Normal} message Normal
                 * @param {$protobuf.IConversionOptions} [options] Conversion options
                 * @returns {Object.<string,*>} Plain object
                 */
                Normal.toObject = function toObject(message, options) {
                    if (!options)
                        options = {};
                    var object = {};
                    if (options.defaults) {
                        object.Pos = null;
                        object.Yaw = 0;
                        object.Duration = 0;
                    }
                    if (message.Pos != null && message.hasOwnProperty("Pos"))
                        object.Pos = $root.message.Point3F.toObject(message.Pos, options);
                    if (message.Yaw != null && message.hasOwnProperty("Yaw"))
                        object.Yaw = options.json && !isFinite(message.Yaw) ? String(message.Yaw) : message.Yaw;
                    if (message.Duration != null && message.hasOwnProperty("Duration"))
                        object.Duration = options.json && !isFinite(message.Duration) ? String(message.Duration) : message.Duration;
                    return object;
                };

                /**
                 * Converts this Normal to JSON.
                 * @function toJSON
                 * @memberof message.C_W_Move.Move.Normal
                 * @instance
                 * @returns {Object.<string,*>} JSON object
                 */
                Normal.prototype.toJSON = function toJSON() {
                    return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
                };

                return Normal;
            })();

            Move.Path = (function() {

                /**
                 * Properties of a Path.
                 * @memberof message.C_W_Move.Move
                 * @interface IPath
                 * @property {number|null} [PathId] Path PathId
                 * @property {number|null} [TimePos] Path TimePos
                 * @property {number|null} [MountId] Path MountId
                 */

                /**
                 * Constructs a new Path.
                 * @memberof message.C_W_Move.Move
                 * @classdesc Represents a Path.
                 * @implements IPath
                 * @constructor
                 * @param {message.C_W_Move.Move.IPath=} [properties] Properties to set
                 */
                function Path(properties) {
                    if (properties)
                        for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                            if (properties[keys[i]] != null)
                                this[keys[i]] = properties[keys[i]];
                }

                /**
                 * Path PathId.
                 * @member {number} PathId
                 * @memberof message.C_W_Move.Move.Path
                 * @instance
                 */
                Path.prototype.PathId = 0;

                /**
                 * Path TimePos.
                 * @member {number} TimePos
                 * @memberof message.C_W_Move.Move.Path
                 * @instance
                 */
                Path.prototype.TimePos = 0;

                /**
                 * Path MountId.
                 * @member {number} MountId
                 * @memberof message.C_W_Move.Move.Path
                 * @instance
                 */
                Path.prototype.MountId = 0;

                /**
                 * Creates a new Path instance using the specified properties.
                 * @function create
                 * @memberof message.C_W_Move.Move.Path
                 * @static
                 * @param {message.C_W_Move.Move.IPath=} [properties] Properties to set
                 * @returns {message.C_W_Move.Move.Path} Path instance
                 */
                Path.create = function create(properties) {
                    return new Path(properties);
                };

                /**
                 * Encodes the specified Path message. Does not implicitly {@link message.C_W_Move.Move.Path.verify|verify} messages.
                 * @function encode
                 * @memberof message.C_W_Move.Move.Path
                 * @static
                 * @param {message.C_W_Move.Move.IPath} message Path message or plain object to encode
                 * @param {$protobuf.Writer} [writer] Writer to encode to
                 * @returns {$protobuf.Writer} Writer
                 */
                Path.encode = function encode(message, writer) {
                    if (!writer)
                        writer = $Writer.create();
                    if (message.PathId != null && message.hasOwnProperty("PathId"))
                        writer.uint32(/* id 1, wireType 0 =*/8).int32(message.PathId);
                    if (message.TimePos != null && message.hasOwnProperty("TimePos"))
                        writer.uint32(/* id 2, wireType 0 =*/16).int32(message.TimePos);
                    if (message.MountId != null && message.hasOwnProperty("MountId"))
                        writer.uint32(/* id 3, wireType 0 =*/24).int32(message.MountId);
                    return writer;
                };

                /**
                 * Encodes the specified Path message, length delimited. Does not implicitly {@link message.C_W_Move.Move.Path.verify|verify} messages.
                 * @function encodeDelimited
                 * @memberof message.C_W_Move.Move.Path
                 * @static
                 * @param {message.C_W_Move.Move.IPath} message Path message or plain object to encode
                 * @param {$protobuf.Writer} [writer] Writer to encode to
                 * @returns {$protobuf.Writer} Writer
                 */
                Path.encodeDelimited = function encodeDelimited(message, writer) {
                    return this.encode(message, writer).ldelim();
                };

                /**
                 * Decodes a Path message from the specified reader or buffer.
                 * @function decode
                 * @memberof message.C_W_Move.Move.Path
                 * @static
                 * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
                 * @param {number} [length] Message length if known beforehand
                 * @returns {message.C_W_Move.Move.Path} Path
                 * @throws {Error} If the payload is not a reader or valid buffer
                 * @throws {$protobuf.util.ProtocolError} If required fields are missing
                 */
                Path.decode = function decode(reader, length) {
                    if (!(reader instanceof $Reader))
                        reader = $Reader.create(reader);
                    var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.C_W_Move.Move.Path();
                    while (reader.pos < end) {
                        var tag = reader.uint32();
                        switch (tag >>> 3) {
                        case 1:
                            message.PathId = reader.int32();
                            break;
                        case 2:
                            message.TimePos = reader.int32();
                            break;
                        case 3:
                            message.MountId = reader.int32();
                            break;
                        default:
                            reader.skipType(tag & 7);
                            break;
                        }
                    }
                    return message;
                };

                /**
                 * Decodes a Path message from the specified reader or buffer, length delimited.
                 * @function decodeDelimited
                 * @memberof message.C_W_Move.Move.Path
                 * @static
                 * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
                 * @returns {message.C_W_Move.Move.Path} Path
                 * @throws {Error} If the payload is not a reader or valid buffer
                 * @throws {$protobuf.util.ProtocolError} If required fields are missing
                 */
                Path.decodeDelimited = function decodeDelimited(reader) {
                    if (!(reader instanceof $Reader))
                        reader = new $Reader(reader);
                    return this.decode(reader, reader.uint32());
                };

                /**
                 * Verifies a Path message.
                 * @function verify
                 * @memberof message.C_W_Move.Move.Path
                 * @static
                 * @param {Object.<string,*>} message Plain object to verify
                 * @returns {string|null} `null` if valid, otherwise the reason why it is not
                 */
                Path.verify = function verify(message) {
                    if (typeof message !== "object" || message === null)
                        return "object expected";
                    if (message.PathId != null && message.hasOwnProperty("PathId"))
                        if (!$util.isInteger(message.PathId))
                            return "PathId: integer expected";
                    if (message.TimePos != null && message.hasOwnProperty("TimePos"))
                        if (!$util.isInteger(message.TimePos))
                            return "TimePos: integer expected";
                    if (message.MountId != null && message.hasOwnProperty("MountId"))
                        if (!$util.isInteger(message.MountId))
                            return "MountId: integer expected";
                    return null;
                };

                /**
                 * Creates a Path message from a plain object. Also converts values to their respective internal types.
                 * @function fromObject
                 * @memberof message.C_W_Move.Move.Path
                 * @static
                 * @param {Object.<string,*>} object Plain object
                 * @returns {message.C_W_Move.Move.Path} Path
                 */
                Path.fromObject = function fromObject(object) {
                    if (object instanceof $root.message.C_W_Move.Move.Path)
                        return object;
                    var message = new $root.message.C_W_Move.Move.Path();
                    if (object.PathId != null)
                        message.PathId = object.PathId | 0;
                    if (object.TimePos != null)
                        message.TimePos = object.TimePos | 0;
                    if (object.MountId != null)
                        message.MountId = object.MountId | 0;
                    return message;
                };

                /**
                 * Creates a plain object from a Path message. Also converts values to other types if specified.
                 * @function toObject
                 * @memberof message.C_W_Move.Move.Path
                 * @static
                 * @param {message.C_W_Move.Move.Path} message Path
                 * @param {$protobuf.IConversionOptions} [options] Conversion options
                 * @returns {Object.<string,*>} Plain object
                 */
                Path.toObject = function toObject(message, options) {
                    if (!options)
                        options = {};
                    var object = {};
                    if (options.defaults) {
                        object.PathId = 0;
                        object.TimePos = 0;
                        object.MountId = 0;
                    }
                    if (message.PathId != null && message.hasOwnProperty("PathId"))
                        object.PathId = message.PathId;
                    if (message.TimePos != null && message.hasOwnProperty("TimePos"))
                        object.TimePos = message.TimePos;
                    if (message.MountId != null && message.hasOwnProperty("MountId"))
                        object.MountId = message.MountId;
                    return object;
                };

                /**
                 * Converts this Path to JSON.
                 * @function toJSON
                 * @memberof message.C_W_Move.Move.Path
                 * @instance
                 * @returns {Object.<string,*>} JSON object
                 */
                Path.prototype.toJSON = function toJSON() {
                    return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
                };

                return Path;
            })();

            Move.Blink = (function() {

                /**
                 * Properties of a Blink.
                 * @memberof message.C_W_Move.Move
                 * @interface IBlink
                 * @property {message.IPoint3F|null} [Pos] Blink Pos
                 * @property {message.IPoint3F|null} [RPos] Blink RPos
                 */

                /**
                 * Constructs a new Blink.
                 * @memberof message.C_W_Move.Move
                 * @classdesc Represents a Blink.
                 * @implements IBlink
                 * @constructor
                 * @param {message.C_W_Move.Move.IBlink=} [properties] Properties to set
                 */
                function Blink(properties) {
                    if (properties)
                        for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                            if (properties[keys[i]] != null)
                                this[keys[i]] = properties[keys[i]];
                }

                /**
                 * Blink Pos.
                 * @member {message.IPoint3F|null|undefined} Pos
                 * @memberof message.C_W_Move.Move.Blink
                 * @instance
                 */
                Blink.prototype.Pos = null;

                /**
                 * Blink RPos.
                 * @member {message.IPoint3F|null|undefined} RPos
                 * @memberof message.C_W_Move.Move.Blink
                 * @instance
                 */
                Blink.prototype.RPos = null;

                /**
                 * Creates a new Blink instance using the specified properties.
                 * @function create
                 * @memberof message.C_W_Move.Move.Blink
                 * @static
                 * @param {message.C_W_Move.Move.IBlink=} [properties] Properties to set
                 * @returns {message.C_W_Move.Move.Blink} Blink instance
                 */
                Blink.create = function create(properties) {
                    return new Blink(properties);
                };

                /**
                 * Encodes the specified Blink message. Does not implicitly {@link message.C_W_Move.Move.Blink.verify|verify} messages.
                 * @function encode
                 * @memberof message.C_W_Move.Move.Blink
                 * @static
                 * @param {message.C_W_Move.Move.IBlink} message Blink message or plain object to encode
                 * @param {$protobuf.Writer} [writer] Writer to encode to
                 * @returns {$protobuf.Writer} Writer
                 */
                Blink.encode = function encode(message, writer) {
                    if (!writer)
                        writer = $Writer.create();
                    if (message.Pos != null && message.hasOwnProperty("Pos"))
                        $root.message.Point3F.encode(message.Pos, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
                    if (message.RPos != null && message.hasOwnProperty("RPos"))
                        $root.message.Point3F.encode(message.RPos, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
                    return writer;
                };

                /**
                 * Encodes the specified Blink message, length delimited. Does not implicitly {@link message.C_W_Move.Move.Blink.verify|verify} messages.
                 * @function encodeDelimited
                 * @memberof message.C_W_Move.Move.Blink
                 * @static
                 * @param {message.C_W_Move.Move.IBlink} message Blink message or plain object to encode
                 * @param {$protobuf.Writer} [writer] Writer to encode to
                 * @returns {$protobuf.Writer} Writer
                 */
                Blink.encodeDelimited = function encodeDelimited(message, writer) {
                    return this.encode(message, writer).ldelim();
                };

                /**
                 * Decodes a Blink message from the specified reader or buffer.
                 * @function decode
                 * @memberof message.C_W_Move.Move.Blink
                 * @static
                 * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
                 * @param {number} [length] Message length if known beforehand
                 * @returns {message.C_W_Move.Move.Blink} Blink
                 * @throws {Error} If the payload is not a reader or valid buffer
                 * @throws {$protobuf.util.ProtocolError} If required fields are missing
                 */
                Blink.decode = function decode(reader, length) {
                    if (!(reader instanceof $Reader))
                        reader = $Reader.create(reader);
                    var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.C_W_Move.Move.Blink();
                    while (reader.pos < end) {
                        var tag = reader.uint32();
                        switch (tag >>> 3) {
                        case 1:
                            message.Pos = $root.message.Point3F.decode(reader, reader.uint32());
                            break;
                        case 2:
                            message.RPos = $root.message.Point3F.decode(reader, reader.uint32());
                            break;
                        default:
                            reader.skipType(tag & 7);
                            break;
                        }
                    }
                    return message;
                };

                /**
                 * Decodes a Blink message from the specified reader or buffer, length delimited.
                 * @function decodeDelimited
                 * @memberof message.C_W_Move.Move.Blink
                 * @static
                 * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
                 * @returns {message.C_W_Move.Move.Blink} Blink
                 * @throws {Error} If the payload is not a reader or valid buffer
                 * @throws {$protobuf.util.ProtocolError} If required fields are missing
                 */
                Blink.decodeDelimited = function decodeDelimited(reader) {
                    if (!(reader instanceof $Reader))
                        reader = new $Reader(reader);
                    return this.decode(reader, reader.uint32());
                };

                /**
                 * Verifies a Blink message.
                 * @function verify
                 * @memberof message.C_W_Move.Move.Blink
                 * @static
                 * @param {Object.<string,*>} message Plain object to verify
                 * @returns {string|null} `null` if valid, otherwise the reason why it is not
                 */
                Blink.verify = function verify(message) {
                    if (typeof message !== "object" || message === null)
                        return "object expected";
                    if (message.Pos != null && message.hasOwnProperty("Pos")) {
                        var error = $root.message.Point3F.verify(message.Pos);
                        if (error)
                            return "Pos." + error;
                    }
                    if (message.RPos != null && message.hasOwnProperty("RPos")) {
                        var error = $root.message.Point3F.verify(message.RPos);
                        if (error)
                            return "RPos." + error;
                    }
                    return null;
                };

                /**
                 * Creates a Blink message from a plain object. Also converts values to their respective internal types.
                 * @function fromObject
                 * @memberof message.C_W_Move.Move.Blink
                 * @static
                 * @param {Object.<string,*>} object Plain object
                 * @returns {message.C_W_Move.Move.Blink} Blink
                 */
                Blink.fromObject = function fromObject(object) {
                    if (object instanceof $root.message.C_W_Move.Move.Blink)
                        return object;
                    var message = new $root.message.C_W_Move.Move.Blink();
                    if (object.Pos != null) {
                        if (typeof object.Pos !== "object")
                            throw TypeError(".message.C_W_Move.Move.Blink.Pos: object expected");
                        message.Pos = $root.message.Point3F.fromObject(object.Pos);
                    }
                    if (object.RPos != null) {
                        if (typeof object.RPos !== "object")
                            throw TypeError(".message.C_W_Move.Move.Blink.RPos: object expected");
                        message.RPos = $root.message.Point3F.fromObject(object.RPos);
                    }
                    return message;
                };

                /**
                 * Creates a plain object from a Blink message. Also converts values to other types if specified.
                 * @function toObject
                 * @memberof message.C_W_Move.Move.Blink
                 * @static
                 * @param {message.C_W_Move.Move.Blink} message Blink
                 * @param {$protobuf.IConversionOptions} [options] Conversion options
                 * @returns {Object.<string,*>} Plain object
                 */
                Blink.toObject = function toObject(message, options) {
                    if (!options)
                        options = {};
                    var object = {};
                    if (options.defaults) {
                        object.Pos = null;
                        object.RPos = null;
                    }
                    if (message.Pos != null && message.hasOwnProperty("Pos"))
                        object.Pos = $root.message.Point3F.toObject(message.Pos, options);
                    if (message.RPos != null && message.hasOwnProperty("RPos"))
                        object.RPos = $root.message.Point3F.toObject(message.RPos, options);
                    return object;
                };

                /**
                 * Converts this Blink to JSON.
                 * @function toJSON
                 * @memberof message.C_W_Move.Move.Blink
                 * @instance
                 * @returns {Object.<string,*>} JSON object
                 */
                Blink.prototype.toJSON = function toJSON() {
                    return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
                };

                return Blink;
            })();

            Move.Jump = (function() {

                /**
                 * Properties of a Jump.
                 * @memberof message.C_W_Move.Move
                 * @interface IJump
                 * @property {message.IPoint3F|null} [BPos] Jump BPos
                 * @property {message.IPoint3F|null} [EPos] Jump EPos
                 * @property {number|null} [Duration] Jump Duration
                 * @property {number|null} [TimePos] Jump TimePos
                 * @property {number|null} [UpExDur] Jump UpExDur
                 * @property {number|null} [DownExDur] Jump DownExDur
                 * @property {number|null} [A] Jump A
                 * @property {number|null} [B] Jump B
                 */

                /**
                 * Constructs a new Jump.
                 * @memberof message.C_W_Move.Move
                 * @classdesc Represents a Jump.
                 * @implements IJump
                 * @constructor
                 * @param {message.C_W_Move.Move.IJump=} [properties] Properties to set
                 */
                function Jump(properties) {
                    if (properties)
                        for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                            if (properties[keys[i]] != null)
                                this[keys[i]] = properties[keys[i]];
                }

                /**
                 * Jump BPos.
                 * @member {message.IPoint3F|null|undefined} BPos
                 * @memberof message.C_W_Move.Move.Jump
                 * @instance
                 */
                Jump.prototype.BPos = null;

                /**
                 * Jump EPos.
                 * @member {message.IPoint3F|null|undefined} EPos
                 * @memberof message.C_W_Move.Move.Jump
                 * @instance
                 */
                Jump.prototype.EPos = null;

                /**
                 * Jump Duration.
                 * @member {number} Duration
                 * @memberof message.C_W_Move.Move.Jump
                 * @instance
                 */
                Jump.prototype.Duration = 0;

                /**
                 * Jump TimePos.
                 * @member {number} TimePos
                 * @memberof message.C_W_Move.Move.Jump
                 * @instance
                 */
                Jump.prototype.TimePos = 0;

                /**
                 * Jump UpExDur.
                 * @member {number} UpExDur
                 * @memberof message.C_W_Move.Move.Jump
                 * @instance
                 */
                Jump.prototype.UpExDur = 0;

                /**
                 * Jump DownExDur.
                 * @member {number} DownExDur
                 * @memberof message.C_W_Move.Move.Jump
                 * @instance
                 */
                Jump.prototype.DownExDur = 0;

                /**
                 * Jump A.
                 * @member {number} A
                 * @memberof message.C_W_Move.Move.Jump
                 * @instance
                 */
                Jump.prototype.A = 0;

                /**
                 * Jump B.
                 * @member {number} B
                 * @memberof message.C_W_Move.Move.Jump
                 * @instance
                 */
                Jump.prototype.B = 0;

                /**
                 * Creates a new Jump instance using the specified properties.
                 * @function create
                 * @memberof message.C_W_Move.Move.Jump
                 * @static
                 * @param {message.C_W_Move.Move.IJump=} [properties] Properties to set
                 * @returns {message.C_W_Move.Move.Jump} Jump instance
                 */
                Jump.create = function create(properties) {
                    return new Jump(properties);
                };

                /**
                 * Encodes the specified Jump message. Does not implicitly {@link message.C_W_Move.Move.Jump.verify|verify} messages.
                 * @function encode
                 * @memberof message.C_W_Move.Move.Jump
                 * @static
                 * @param {message.C_W_Move.Move.IJump} message Jump message or plain object to encode
                 * @param {$protobuf.Writer} [writer] Writer to encode to
                 * @returns {$protobuf.Writer} Writer
                 */
                Jump.encode = function encode(message, writer) {
                    if (!writer)
                        writer = $Writer.create();
                    if (message.BPos != null && message.hasOwnProperty("BPos"))
                        $root.message.Point3F.encode(message.BPos, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
                    if (message.EPos != null && message.hasOwnProperty("EPos"))
                        $root.message.Point3F.encode(message.EPos, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
                    if (message.Duration != null && message.hasOwnProperty("Duration"))
                        writer.uint32(/* id 3, wireType 0 =*/24).int32(message.Duration);
                    if (message.TimePos != null && message.hasOwnProperty("TimePos"))
                        writer.uint32(/* id 4, wireType 0 =*/32).int32(message.TimePos);
                    if (message.UpExDur != null && message.hasOwnProperty("UpExDur"))
                        writer.uint32(/* id 5, wireType 0 =*/40).int32(message.UpExDur);
                    if (message.DownExDur != null && message.hasOwnProperty("DownExDur"))
                        writer.uint32(/* id 6, wireType 0 =*/48).int32(message.DownExDur);
                    if (message.A != null && message.hasOwnProperty("A"))
                        writer.uint32(/* id 7, wireType 0 =*/56).int32(message.A);
                    if (message.B != null && message.hasOwnProperty("B"))
                        writer.uint32(/* id 8, wireType 0 =*/64).int32(message.B);
                    return writer;
                };

                /**
                 * Encodes the specified Jump message, length delimited. Does not implicitly {@link message.C_W_Move.Move.Jump.verify|verify} messages.
                 * @function encodeDelimited
                 * @memberof message.C_W_Move.Move.Jump
                 * @static
                 * @param {message.C_W_Move.Move.IJump} message Jump message or plain object to encode
                 * @param {$protobuf.Writer} [writer] Writer to encode to
                 * @returns {$protobuf.Writer} Writer
                 */
                Jump.encodeDelimited = function encodeDelimited(message, writer) {
                    return this.encode(message, writer).ldelim();
                };

                /**
                 * Decodes a Jump message from the specified reader or buffer.
                 * @function decode
                 * @memberof message.C_W_Move.Move.Jump
                 * @static
                 * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
                 * @param {number} [length] Message length if known beforehand
                 * @returns {message.C_W_Move.Move.Jump} Jump
                 * @throws {Error} If the payload is not a reader or valid buffer
                 * @throws {$protobuf.util.ProtocolError} If required fields are missing
                 */
                Jump.decode = function decode(reader, length) {
                    if (!(reader instanceof $Reader))
                        reader = $Reader.create(reader);
                    var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.C_W_Move.Move.Jump();
                    while (reader.pos < end) {
                        var tag = reader.uint32();
                        switch (tag >>> 3) {
                        case 1:
                            message.BPos = $root.message.Point3F.decode(reader, reader.uint32());
                            break;
                        case 2:
                            message.EPos = $root.message.Point3F.decode(reader, reader.uint32());
                            break;
                        case 3:
                            message.Duration = reader.int32();
                            break;
                        case 4:
                            message.TimePos = reader.int32();
                            break;
                        case 5:
                            message.UpExDur = reader.int32();
                            break;
                        case 6:
                            message.DownExDur = reader.int32();
                            break;
                        case 7:
                            message.A = reader.int32();
                            break;
                        case 8:
                            message.B = reader.int32();
                            break;
                        default:
                            reader.skipType(tag & 7);
                            break;
                        }
                    }
                    return message;
                };

                /**
                 * Decodes a Jump message from the specified reader or buffer, length delimited.
                 * @function decodeDelimited
                 * @memberof message.C_W_Move.Move.Jump
                 * @static
                 * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
                 * @returns {message.C_W_Move.Move.Jump} Jump
                 * @throws {Error} If the payload is not a reader or valid buffer
                 * @throws {$protobuf.util.ProtocolError} If required fields are missing
                 */
                Jump.decodeDelimited = function decodeDelimited(reader) {
                    if (!(reader instanceof $Reader))
                        reader = new $Reader(reader);
                    return this.decode(reader, reader.uint32());
                };

                /**
                 * Verifies a Jump message.
                 * @function verify
                 * @memberof message.C_W_Move.Move.Jump
                 * @static
                 * @param {Object.<string,*>} message Plain object to verify
                 * @returns {string|null} `null` if valid, otherwise the reason why it is not
                 */
                Jump.verify = function verify(message) {
                    if (typeof message !== "object" || message === null)
                        return "object expected";
                    if (message.BPos != null && message.hasOwnProperty("BPos")) {
                        var error = $root.message.Point3F.verify(message.BPos);
                        if (error)
                            return "BPos." + error;
                    }
                    if (message.EPos != null && message.hasOwnProperty("EPos")) {
                        var error = $root.message.Point3F.verify(message.EPos);
                        if (error)
                            return "EPos." + error;
                    }
                    if (message.Duration != null && message.hasOwnProperty("Duration"))
                        if (!$util.isInteger(message.Duration))
                            return "Duration: integer expected";
                    if (message.TimePos != null && message.hasOwnProperty("TimePos"))
                        if (!$util.isInteger(message.TimePos))
                            return "TimePos: integer expected";
                    if (message.UpExDur != null && message.hasOwnProperty("UpExDur"))
                        if (!$util.isInteger(message.UpExDur))
                            return "UpExDur: integer expected";
                    if (message.DownExDur != null && message.hasOwnProperty("DownExDur"))
                        if (!$util.isInteger(message.DownExDur))
                            return "DownExDur: integer expected";
                    if (message.A != null && message.hasOwnProperty("A"))
                        if (!$util.isInteger(message.A))
                            return "A: integer expected";
                    if (message.B != null && message.hasOwnProperty("B"))
                        if (!$util.isInteger(message.B))
                            return "B: integer expected";
                    return null;
                };

                /**
                 * Creates a Jump message from a plain object. Also converts values to their respective internal types.
                 * @function fromObject
                 * @memberof message.C_W_Move.Move.Jump
                 * @static
                 * @param {Object.<string,*>} object Plain object
                 * @returns {message.C_W_Move.Move.Jump} Jump
                 */
                Jump.fromObject = function fromObject(object) {
                    if (object instanceof $root.message.C_W_Move.Move.Jump)
                        return object;
                    var message = new $root.message.C_W_Move.Move.Jump();
                    if (object.BPos != null) {
                        if (typeof object.BPos !== "object")
                            throw TypeError(".message.C_W_Move.Move.Jump.BPos: object expected");
                        message.BPos = $root.message.Point3F.fromObject(object.BPos);
                    }
                    if (object.EPos != null) {
                        if (typeof object.EPos !== "object")
                            throw TypeError(".message.C_W_Move.Move.Jump.EPos: object expected");
                        message.EPos = $root.message.Point3F.fromObject(object.EPos);
                    }
                    if (object.Duration != null)
                        message.Duration = object.Duration | 0;
                    if (object.TimePos != null)
                        message.TimePos = object.TimePos | 0;
                    if (object.UpExDur != null)
                        message.UpExDur = object.UpExDur | 0;
                    if (object.DownExDur != null)
                        message.DownExDur = object.DownExDur | 0;
                    if (object.A != null)
                        message.A = object.A | 0;
                    if (object.B != null)
                        message.B = object.B | 0;
                    return message;
                };

                /**
                 * Creates a plain object from a Jump message. Also converts values to other types if specified.
                 * @function toObject
                 * @memberof message.C_W_Move.Move.Jump
                 * @static
                 * @param {message.C_W_Move.Move.Jump} message Jump
                 * @param {$protobuf.IConversionOptions} [options] Conversion options
                 * @returns {Object.<string,*>} Plain object
                 */
                Jump.toObject = function toObject(message, options) {
                    if (!options)
                        options = {};
                    var object = {};
                    if (options.defaults) {
                        object.BPos = null;
                        object.EPos = null;
                        object.Duration = 0;
                        object.TimePos = 0;
                        object.UpExDur = 0;
                        object.DownExDur = 0;
                        object.A = 0;
                        object.B = 0;
                    }
                    if (message.BPos != null && message.hasOwnProperty("BPos"))
                        object.BPos = $root.message.Point3F.toObject(message.BPos, options);
                    if (message.EPos != null && message.hasOwnProperty("EPos"))
                        object.EPos = $root.message.Point3F.toObject(message.EPos, options);
                    if (message.Duration != null && message.hasOwnProperty("Duration"))
                        object.Duration = message.Duration;
                    if (message.TimePos != null && message.hasOwnProperty("TimePos"))
                        object.TimePos = message.TimePos;
                    if (message.UpExDur != null && message.hasOwnProperty("UpExDur"))
                        object.UpExDur = message.UpExDur;
                    if (message.DownExDur != null && message.hasOwnProperty("DownExDur"))
                        object.DownExDur = message.DownExDur;
                    if (message.A != null && message.hasOwnProperty("A"))
                        object.A = message.A;
                    if (message.B != null && message.hasOwnProperty("B"))
                        object.B = message.B;
                    return object;
                };

                /**
                 * Converts this Jump to JSON.
                 * @function toJSON
                 * @memberof message.C_W_Move.Move.Jump
                 * @instance
                 * @returns {Object.<string,*>} JSON object
                 */
                Jump.prototype.toJSON = function toJSON() {
                    return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
                };

                return Jump;
            })();

            Move.Line = (function() {

                /**
                 * Properties of a Line.
                 * @memberof message.C_W_Move.Move
                 * @interface ILine
                 * @property {message.IPoint3F|null} [BPos] Line BPos
                 * @property {message.IPoint3F|null} [EPos] Line EPos
                 * @property {number|null} [Duration] Line Duration
                 * @property {number|null} [TimePos] Line TimePos
                 */

                /**
                 * Constructs a new Line.
                 * @memberof message.C_W_Move.Move
                 * @classdesc Represents a Line.
                 * @implements ILine
                 * @constructor
                 * @param {message.C_W_Move.Move.ILine=} [properties] Properties to set
                 */
                function Line(properties) {
                    if (properties)
                        for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                            if (properties[keys[i]] != null)
                                this[keys[i]] = properties[keys[i]];
                }

                /**
                 * Line BPos.
                 * @member {message.IPoint3F|null|undefined} BPos
                 * @memberof message.C_W_Move.Move.Line
                 * @instance
                 */
                Line.prototype.BPos = null;

                /**
                 * Line EPos.
                 * @member {message.IPoint3F|null|undefined} EPos
                 * @memberof message.C_W_Move.Move.Line
                 * @instance
                 */
                Line.prototype.EPos = null;

                /**
                 * Line Duration.
                 * @member {number} Duration
                 * @memberof message.C_W_Move.Move.Line
                 * @instance
                 */
                Line.prototype.Duration = 0;

                /**
                 * Line TimePos.
                 * @member {number} TimePos
                 * @memberof message.C_W_Move.Move.Line
                 * @instance
                 */
                Line.prototype.TimePos = 0;

                /**
                 * Creates a new Line instance using the specified properties.
                 * @function create
                 * @memberof message.C_W_Move.Move.Line
                 * @static
                 * @param {message.C_W_Move.Move.ILine=} [properties] Properties to set
                 * @returns {message.C_W_Move.Move.Line} Line instance
                 */
                Line.create = function create(properties) {
                    return new Line(properties);
                };

                /**
                 * Encodes the specified Line message. Does not implicitly {@link message.C_W_Move.Move.Line.verify|verify} messages.
                 * @function encode
                 * @memberof message.C_W_Move.Move.Line
                 * @static
                 * @param {message.C_W_Move.Move.ILine} message Line message or plain object to encode
                 * @param {$protobuf.Writer} [writer] Writer to encode to
                 * @returns {$protobuf.Writer} Writer
                 */
                Line.encode = function encode(message, writer) {
                    if (!writer)
                        writer = $Writer.create();
                    if (message.BPos != null && message.hasOwnProperty("BPos"))
                        $root.message.Point3F.encode(message.BPos, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
                    if (message.EPos != null && message.hasOwnProperty("EPos"))
                        $root.message.Point3F.encode(message.EPos, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
                    if (message.Duration != null && message.hasOwnProperty("Duration"))
                        writer.uint32(/* id 3, wireType 0 =*/24).int32(message.Duration);
                    if (message.TimePos != null && message.hasOwnProperty("TimePos"))
                        writer.uint32(/* id 4, wireType 0 =*/32).int32(message.TimePos);
                    return writer;
                };

                /**
                 * Encodes the specified Line message, length delimited. Does not implicitly {@link message.C_W_Move.Move.Line.verify|verify} messages.
                 * @function encodeDelimited
                 * @memberof message.C_W_Move.Move.Line
                 * @static
                 * @param {message.C_W_Move.Move.ILine} message Line message or plain object to encode
                 * @param {$protobuf.Writer} [writer] Writer to encode to
                 * @returns {$protobuf.Writer} Writer
                 */
                Line.encodeDelimited = function encodeDelimited(message, writer) {
                    return this.encode(message, writer).ldelim();
                };

                /**
                 * Decodes a Line message from the specified reader or buffer.
                 * @function decode
                 * @memberof message.C_W_Move.Move.Line
                 * @static
                 * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
                 * @param {number} [length] Message length if known beforehand
                 * @returns {message.C_W_Move.Move.Line} Line
                 * @throws {Error} If the payload is not a reader or valid buffer
                 * @throws {$protobuf.util.ProtocolError} If required fields are missing
                 */
                Line.decode = function decode(reader, length) {
                    if (!(reader instanceof $Reader))
                        reader = $Reader.create(reader);
                    var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.C_W_Move.Move.Line();
                    while (reader.pos < end) {
                        var tag = reader.uint32();
                        switch (tag >>> 3) {
                        case 1:
                            message.BPos = $root.message.Point3F.decode(reader, reader.uint32());
                            break;
                        case 2:
                            message.EPos = $root.message.Point3F.decode(reader, reader.uint32());
                            break;
                        case 3:
                            message.Duration = reader.int32();
                            break;
                        case 4:
                            message.TimePos = reader.int32();
                            break;
                        default:
                            reader.skipType(tag & 7);
                            break;
                        }
                    }
                    return message;
                };

                /**
                 * Decodes a Line message from the specified reader or buffer, length delimited.
                 * @function decodeDelimited
                 * @memberof message.C_W_Move.Move.Line
                 * @static
                 * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
                 * @returns {message.C_W_Move.Move.Line} Line
                 * @throws {Error} If the payload is not a reader or valid buffer
                 * @throws {$protobuf.util.ProtocolError} If required fields are missing
                 */
                Line.decodeDelimited = function decodeDelimited(reader) {
                    if (!(reader instanceof $Reader))
                        reader = new $Reader(reader);
                    return this.decode(reader, reader.uint32());
                };

                /**
                 * Verifies a Line message.
                 * @function verify
                 * @memberof message.C_W_Move.Move.Line
                 * @static
                 * @param {Object.<string,*>} message Plain object to verify
                 * @returns {string|null} `null` if valid, otherwise the reason why it is not
                 */
                Line.verify = function verify(message) {
                    if (typeof message !== "object" || message === null)
                        return "object expected";
                    if (message.BPos != null && message.hasOwnProperty("BPos")) {
                        var error = $root.message.Point3F.verify(message.BPos);
                        if (error)
                            return "BPos." + error;
                    }
                    if (message.EPos != null && message.hasOwnProperty("EPos")) {
                        var error = $root.message.Point3F.verify(message.EPos);
                        if (error)
                            return "EPos." + error;
                    }
                    if (message.Duration != null && message.hasOwnProperty("Duration"))
                        if (!$util.isInteger(message.Duration))
                            return "Duration: integer expected";
                    if (message.TimePos != null && message.hasOwnProperty("TimePos"))
                        if (!$util.isInteger(message.TimePos))
                            return "TimePos: integer expected";
                    return null;
                };

                /**
                 * Creates a Line message from a plain object. Also converts values to their respective internal types.
                 * @function fromObject
                 * @memberof message.C_W_Move.Move.Line
                 * @static
                 * @param {Object.<string,*>} object Plain object
                 * @returns {message.C_W_Move.Move.Line} Line
                 */
                Line.fromObject = function fromObject(object) {
                    if (object instanceof $root.message.C_W_Move.Move.Line)
                        return object;
                    var message = new $root.message.C_W_Move.Move.Line();
                    if (object.BPos != null) {
                        if (typeof object.BPos !== "object")
                            throw TypeError(".message.C_W_Move.Move.Line.BPos: object expected");
                        message.BPos = $root.message.Point3F.fromObject(object.BPos);
                    }
                    if (object.EPos != null) {
                        if (typeof object.EPos !== "object")
                            throw TypeError(".message.C_W_Move.Move.Line.EPos: object expected");
                        message.EPos = $root.message.Point3F.fromObject(object.EPos);
                    }
                    if (object.Duration != null)
                        message.Duration = object.Duration | 0;
                    if (object.TimePos != null)
                        message.TimePos = object.TimePos | 0;
                    return message;
                };

                /**
                 * Creates a plain object from a Line message. Also converts values to other types if specified.
                 * @function toObject
                 * @memberof message.C_W_Move.Move.Line
                 * @static
                 * @param {message.C_W_Move.Move.Line} message Line
                 * @param {$protobuf.IConversionOptions} [options] Conversion options
                 * @returns {Object.<string,*>} Plain object
                 */
                Line.toObject = function toObject(message, options) {
                    if (!options)
                        options = {};
                    var object = {};
                    if (options.defaults) {
                        object.BPos = null;
                        object.EPos = null;
                        object.Duration = 0;
                        object.TimePos = 0;
                    }
                    if (message.BPos != null && message.hasOwnProperty("BPos"))
                        object.BPos = $root.message.Point3F.toObject(message.BPos, options);
                    if (message.EPos != null && message.hasOwnProperty("EPos"))
                        object.EPos = $root.message.Point3F.toObject(message.EPos, options);
                    if (message.Duration != null && message.hasOwnProperty("Duration"))
                        object.Duration = message.Duration;
                    if (message.TimePos != null && message.hasOwnProperty("TimePos"))
                        object.TimePos = message.TimePos;
                    return object;
                };

                /**
                 * Converts this Line to JSON.
                 * @function toJSON
                 * @memberof message.C_W_Move.Move.Line
                 * @instance
                 * @returns {Object.<string,*>} JSON object
                 */
                Line.prototype.toJSON = function toJSON() {
                    return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
                };

                return Line;
            })();

            return Move;
        })();

        return C_W_Move;
    })();

    message.W_C_LoginMap = (function() {

        /**
         * Properties of a W_C_LoginMap.
         * @memberof message
         * @interface IW_C_LoginMap
         * @property {message.IIpacket|null} [PacketHead] W_C_LoginMap PacketHead
         * @property {number|Long|null} [Id] W_C_LoginMap Id
         * @property {message.IPoint3F|null} [Pos] W_C_LoginMap Pos
         * @property {number|null} [Rotation] W_C_LoginMap Rotation
         */

        /**
         * Constructs a new W_C_LoginMap.
         * @memberof message
         * @classdesc Represents a W_C_LoginMap.
         * @implements IW_C_LoginMap
         * @constructor
         * @param {message.IW_C_LoginMap=} [properties] Properties to set
         */
        function W_C_LoginMap(properties) {
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * W_C_LoginMap PacketHead.
         * @member {message.IIpacket|null|undefined} PacketHead
         * @memberof message.W_C_LoginMap
         * @instance
         */
        W_C_LoginMap.prototype.PacketHead = null;

        /**
         * W_C_LoginMap Id.
         * @member {number|Long} Id
         * @memberof message.W_C_LoginMap
         * @instance
         */
        W_C_LoginMap.prototype.Id = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * W_C_LoginMap Pos.
         * @member {message.IPoint3F|null|undefined} Pos
         * @memberof message.W_C_LoginMap
         * @instance
         */
        W_C_LoginMap.prototype.Pos = null;

        /**
         * W_C_LoginMap Rotation.
         * @member {number} Rotation
         * @memberof message.W_C_LoginMap
         * @instance
         */
        W_C_LoginMap.prototype.Rotation = 0;

        /**
         * Creates a new W_C_LoginMap instance using the specified properties.
         * @function create
         * @memberof message.W_C_LoginMap
         * @static
         * @param {message.IW_C_LoginMap=} [properties] Properties to set
         * @returns {message.W_C_LoginMap} W_C_LoginMap instance
         */
        W_C_LoginMap.create = function create(properties) {
            return new W_C_LoginMap(properties);
        };

        /**
         * Encodes the specified W_C_LoginMap message. Does not implicitly {@link message.W_C_LoginMap.verify|verify} messages.
         * @function encode
         * @memberof message.W_C_LoginMap
         * @static
         * @param {message.IW_C_LoginMap} message W_C_LoginMap message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        W_C_LoginMap.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                $root.message.Ipacket.encode(message.PacketHead, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
            if (message.Id != null && message.hasOwnProperty("Id"))
                writer.uint32(/* id 2, wireType 0 =*/16).int64(message.Id);
            if (message.Pos != null && message.hasOwnProperty("Pos"))
                $root.message.Point3F.encode(message.Pos, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
            if (message.Rotation != null && message.hasOwnProperty("Rotation"))
                writer.uint32(/* id 4, wireType 5 =*/37).float(message.Rotation);
            return writer;
        };

        /**
         * Encodes the specified W_C_LoginMap message, length delimited. Does not implicitly {@link message.W_C_LoginMap.verify|verify} messages.
         * @function encodeDelimited
         * @memberof message.W_C_LoginMap
         * @static
         * @param {message.IW_C_LoginMap} message W_C_LoginMap message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        W_C_LoginMap.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a W_C_LoginMap message from the specified reader or buffer.
         * @function decode
         * @memberof message.W_C_LoginMap
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {message.W_C_LoginMap} W_C_LoginMap
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        W_C_LoginMap.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.W_C_LoginMap();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.PacketHead = $root.message.Ipacket.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.Id = reader.int64();
                    break;
                case 3:
                    message.Pos = $root.message.Point3F.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.Rotation = reader.float();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a W_C_LoginMap message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof message.W_C_LoginMap
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {message.W_C_LoginMap} W_C_LoginMap
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        W_C_LoginMap.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a W_C_LoginMap message.
         * @function verify
         * @memberof message.W_C_LoginMap
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        W_C_LoginMap.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead")) {
                var error = $root.message.Ipacket.verify(message.PacketHead);
                if (error)
                    return "PacketHead." + error;
            }
            if (message.Id != null && message.hasOwnProperty("Id"))
                if (!$util.isInteger(message.Id) && !(message.Id && $util.isInteger(message.Id.low) && $util.isInteger(message.Id.high)))
                    return "Id: integer|Long expected";
            if (message.Pos != null && message.hasOwnProperty("Pos")) {
                var error = $root.message.Point3F.verify(message.Pos);
                if (error)
                    return "Pos." + error;
            }
            if (message.Rotation != null && message.hasOwnProperty("Rotation"))
                if (typeof message.Rotation !== "number")
                    return "Rotation: number expected";
            return null;
        };

        /**
         * Creates a W_C_LoginMap message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof message.W_C_LoginMap
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {message.W_C_LoginMap} W_C_LoginMap
         */
        W_C_LoginMap.fromObject = function fromObject(object) {
            if (object instanceof $root.message.W_C_LoginMap)
                return object;
            var message = new $root.message.W_C_LoginMap();
            if (object.PacketHead != null) {
                if (typeof object.PacketHead !== "object")
                    throw TypeError(".message.W_C_LoginMap.PacketHead: object expected");
                message.PacketHead = $root.message.Ipacket.fromObject(object.PacketHead);
            }
            if (object.Id != null)
                if ($util.Long)
                    (message.Id = $util.Long.fromValue(object.Id)).unsigned = false;
                else if (typeof object.Id === "string")
                    message.Id = parseInt(object.Id, 10);
                else if (typeof object.Id === "number")
                    message.Id = object.Id;
                else if (typeof object.Id === "object")
                    message.Id = new $util.LongBits(object.Id.low >>> 0, object.Id.high >>> 0).toNumber();
            if (object.Pos != null) {
                if (typeof object.Pos !== "object")
                    throw TypeError(".message.W_C_LoginMap.Pos: object expected");
                message.Pos = $root.message.Point3F.fromObject(object.Pos);
            }
            if (object.Rotation != null)
                message.Rotation = Number(object.Rotation);
            return message;
        };

        /**
         * Creates a plain object from a W_C_LoginMap message. Also converts values to other types if specified.
         * @function toObject
         * @memberof message.W_C_LoginMap
         * @static
         * @param {message.W_C_LoginMap} message W_C_LoginMap
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        W_C_LoginMap.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.defaults) {
                object.PacketHead = null;
                if ($util.Long) {
                    var long = new $util.Long(0, 0, false);
                    object.Id = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.Id = options.longs === String ? "0" : 0;
                object.Pos = null;
                object.Rotation = 0;
            }
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                object.PacketHead = $root.message.Ipacket.toObject(message.PacketHead, options);
            if (message.Id != null && message.hasOwnProperty("Id"))
                if (typeof message.Id === "number")
                    object.Id = options.longs === String ? String(message.Id) : message.Id;
                else
                    object.Id = options.longs === String ? $util.Long.prototype.toString.call(message.Id) : options.longs === Number ? new $util.LongBits(message.Id.low >>> 0, message.Id.high >>> 0).toNumber() : message.Id;
            if (message.Pos != null && message.hasOwnProperty("Pos"))
                object.Pos = $root.message.Point3F.toObject(message.Pos, options);
            if (message.Rotation != null && message.hasOwnProperty("Rotation"))
                object.Rotation = options.json && !isFinite(message.Rotation) ? String(message.Rotation) : message.Rotation;
            return object;
        };

        /**
         * Converts this W_C_LoginMap to JSON.
         * @function toJSON
         * @memberof message.W_C_LoginMap
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        W_C_LoginMap.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return W_C_LoginMap;
    })();

    message.W_C_Move = (function() {

        /**
         * Properties of a W_C_Move.
         * @memberof message
         * @interface IW_C_Move
         * @property {message.IIpacket|null} [PacketHead] W_C_Move PacketHead
         * @property {number|Long|null} [Id] W_C_Move Id
         * @property {message.IPoint3F|null} [Pos] W_C_Move Pos
         * @property {number|null} [Rotation] W_C_Move Rotation
         */

        /**
         * Constructs a new W_C_Move.
         * @memberof message
         * @classdesc Represents a W_C_Move.
         * @implements IW_C_Move
         * @constructor
         * @param {message.IW_C_Move=} [properties] Properties to set
         */
        function W_C_Move(properties) {
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * W_C_Move PacketHead.
         * @member {message.IIpacket|null|undefined} PacketHead
         * @memberof message.W_C_Move
         * @instance
         */
        W_C_Move.prototype.PacketHead = null;

        /**
         * W_C_Move Id.
         * @member {number|Long} Id
         * @memberof message.W_C_Move
         * @instance
         */
        W_C_Move.prototype.Id = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * W_C_Move Pos.
         * @member {message.IPoint3F|null|undefined} Pos
         * @memberof message.W_C_Move
         * @instance
         */
        W_C_Move.prototype.Pos = null;

        /**
         * W_C_Move Rotation.
         * @member {number} Rotation
         * @memberof message.W_C_Move
         * @instance
         */
        W_C_Move.prototype.Rotation = 0;

        /**
         * Creates a new W_C_Move instance using the specified properties.
         * @function create
         * @memberof message.W_C_Move
         * @static
         * @param {message.IW_C_Move=} [properties] Properties to set
         * @returns {message.W_C_Move} W_C_Move instance
         */
        W_C_Move.create = function create(properties) {
            return new W_C_Move(properties);
        };

        /**
         * Encodes the specified W_C_Move message. Does not implicitly {@link message.W_C_Move.verify|verify} messages.
         * @function encode
         * @memberof message.W_C_Move
         * @static
         * @param {message.IW_C_Move} message W_C_Move message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        W_C_Move.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                $root.message.Ipacket.encode(message.PacketHead, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
            if (message.Id != null && message.hasOwnProperty("Id"))
                writer.uint32(/* id 2, wireType 0 =*/16).int64(message.Id);
            if (message.Pos != null && message.hasOwnProperty("Pos"))
                $root.message.Point3F.encode(message.Pos, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
            if (message.Rotation != null && message.hasOwnProperty("Rotation"))
                writer.uint32(/* id 4, wireType 5 =*/37).float(message.Rotation);
            return writer;
        };

        /**
         * Encodes the specified W_C_Move message, length delimited. Does not implicitly {@link message.W_C_Move.verify|verify} messages.
         * @function encodeDelimited
         * @memberof message.W_C_Move
         * @static
         * @param {message.IW_C_Move} message W_C_Move message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        W_C_Move.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a W_C_Move message from the specified reader or buffer.
         * @function decode
         * @memberof message.W_C_Move
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {message.W_C_Move} W_C_Move
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        W_C_Move.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.W_C_Move();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.PacketHead = $root.message.Ipacket.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.Id = reader.int64();
                    break;
                case 3:
                    message.Pos = $root.message.Point3F.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.Rotation = reader.float();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a W_C_Move message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof message.W_C_Move
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {message.W_C_Move} W_C_Move
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        W_C_Move.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a W_C_Move message.
         * @function verify
         * @memberof message.W_C_Move
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        W_C_Move.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead")) {
                var error = $root.message.Ipacket.verify(message.PacketHead);
                if (error)
                    return "PacketHead." + error;
            }
            if (message.Id != null && message.hasOwnProperty("Id"))
                if (!$util.isInteger(message.Id) && !(message.Id && $util.isInteger(message.Id.low) && $util.isInteger(message.Id.high)))
                    return "Id: integer|Long expected";
            if (message.Pos != null && message.hasOwnProperty("Pos")) {
                var error = $root.message.Point3F.verify(message.Pos);
                if (error)
                    return "Pos." + error;
            }
            if (message.Rotation != null && message.hasOwnProperty("Rotation"))
                if (typeof message.Rotation !== "number")
                    return "Rotation: number expected";
            return null;
        };

        /**
         * Creates a W_C_Move message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof message.W_C_Move
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {message.W_C_Move} W_C_Move
         */
        W_C_Move.fromObject = function fromObject(object) {
            if (object instanceof $root.message.W_C_Move)
                return object;
            var message = new $root.message.W_C_Move();
            if (object.PacketHead != null) {
                if (typeof object.PacketHead !== "object")
                    throw TypeError(".message.W_C_Move.PacketHead: object expected");
                message.PacketHead = $root.message.Ipacket.fromObject(object.PacketHead);
            }
            if (object.Id != null)
                if ($util.Long)
                    (message.Id = $util.Long.fromValue(object.Id)).unsigned = false;
                else if (typeof object.Id === "string")
                    message.Id = parseInt(object.Id, 10);
                else if (typeof object.Id === "number")
                    message.Id = object.Id;
                else if (typeof object.Id === "object")
                    message.Id = new $util.LongBits(object.Id.low >>> 0, object.Id.high >>> 0).toNumber();
            if (object.Pos != null) {
                if (typeof object.Pos !== "object")
                    throw TypeError(".message.W_C_Move.Pos: object expected");
                message.Pos = $root.message.Point3F.fromObject(object.Pos);
            }
            if (object.Rotation != null)
                message.Rotation = Number(object.Rotation);
            return message;
        };

        /**
         * Creates a plain object from a W_C_Move message. Also converts values to other types if specified.
         * @function toObject
         * @memberof message.W_C_Move
         * @static
         * @param {message.W_C_Move} message W_C_Move
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        W_C_Move.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.defaults) {
                object.PacketHead = null;
                if ($util.Long) {
                    var long = new $util.Long(0, 0, false);
                    object.Id = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.Id = options.longs === String ? "0" : 0;
                object.Pos = null;
                object.Rotation = 0;
            }
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                object.PacketHead = $root.message.Ipacket.toObject(message.PacketHead, options);
            if (message.Id != null && message.hasOwnProperty("Id"))
                if (typeof message.Id === "number")
                    object.Id = options.longs === String ? String(message.Id) : message.Id;
                else
                    object.Id = options.longs === String ? $util.Long.prototype.toString.call(message.Id) : options.longs === Number ? new $util.LongBits(message.Id.low >>> 0, message.Id.high >>> 0).toNumber() : message.Id;
            if (message.Pos != null && message.hasOwnProperty("Pos"))
                object.Pos = $root.message.Point3F.toObject(message.Pos, options);
            if (message.Rotation != null && message.hasOwnProperty("Rotation"))
                object.Rotation = options.json && !isFinite(message.Rotation) ? String(message.Rotation) : message.Rotation;
            return object;
        };

        /**
         * Converts this W_C_Move to JSON.
         * @function toJSON
         * @memberof message.W_C_Move
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        W_C_Move.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return W_C_Move;
    })();

    message.W_C_ADD_SIMOBJ = (function() {

        /**
         * Properties of a W_C_ADD_SIMOBJ.
         * @memberof message
         * @interface IW_C_ADD_SIMOBJ
         * @property {message.IIpacket|null} [PacketHead] W_C_ADD_SIMOBJ PacketHead
         * @property {number|Long|null} [Id] W_C_ADD_SIMOBJ Id
         * @property {message.IPoint3F|null} [Pos] W_C_ADD_SIMOBJ Pos
         * @property {number|null} [Rotation] W_C_ADD_SIMOBJ Rotation
         */

        /**
         * Constructs a new W_C_ADD_SIMOBJ.
         * @memberof message
         * @classdesc Represents a W_C_ADD_SIMOBJ.
         * @implements IW_C_ADD_SIMOBJ
         * @constructor
         * @param {message.IW_C_ADD_SIMOBJ=} [properties] Properties to set
         */
        function W_C_ADD_SIMOBJ(properties) {
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * W_C_ADD_SIMOBJ PacketHead.
         * @member {message.IIpacket|null|undefined} PacketHead
         * @memberof message.W_C_ADD_SIMOBJ
         * @instance
         */
        W_C_ADD_SIMOBJ.prototype.PacketHead = null;

        /**
         * W_C_ADD_SIMOBJ Id.
         * @member {number|Long} Id
         * @memberof message.W_C_ADD_SIMOBJ
         * @instance
         */
        W_C_ADD_SIMOBJ.prototype.Id = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * W_C_ADD_SIMOBJ Pos.
         * @member {message.IPoint3F|null|undefined} Pos
         * @memberof message.W_C_ADD_SIMOBJ
         * @instance
         */
        W_C_ADD_SIMOBJ.prototype.Pos = null;

        /**
         * W_C_ADD_SIMOBJ Rotation.
         * @member {number} Rotation
         * @memberof message.W_C_ADD_SIMOBJ
         * @instance
         */
        W_C_ADD_SIMOBJ.prototype.Rotation = 0;

        /**
         * Creates a new W_C_ADD_SIMOBJ instance using the specified properties.
         * @function create
         * @memberof message.W_C_ADD_SIMOBJ
         * @static
         * @param {message.IW_C_ADD_SIMOBJ=} [properties] Properties to set
         * @returns {message.W_C_ADD_SIMOBJ} W_C_ADD_SIMOBJ instance
         */
        W_C_ADD_SIMOBJ.create = function create(properties) {
            return new W_C_ADD_SIMOBJ(properties);
        };

        /**
         * Encodes the specified W_C_ADD_SIMOBJ message. Does not implicitly {@link message.W_C_ADD_SIMOBJ.verify|verify} messages.
         * @function encode
         * @memberof message.W_C_ADD_SIMOBJ
         * @static
         * @param {message.IW_C_ADD_SIMOBJ} message W_C_ADD_SIMOBJ message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        W_C_ADD_SIMOBJ.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                $root.message.Ipacket.encode(message.PacketHead, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
            if (message.Id != null && message.hasOwnProperty("Id"))
                writer.uint32(/* id 2, wireType 0 =*/16).int64(message.Id);
            if (message.Pos != null && message.hasOwnProperty("Pos"))
                $root.message.Point3F.encode(message.Pos, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
            if (message.Rotation != null && message.hasOwnProperty("Rotation"))
                writer.uint32(/* id 4, wireType 5 =*/37).float(message.Rotation);
            return writer;
        };

        /**
         * Encodes the specified W_C_ADD_SIMOBJ message, length delimited. Does not implicitly {@link message.W_C_ADD_SIMOBJ.verify|verify} messages.
         * @function encodeDelimited
         * @memberof message.W_C_ADD_SIMOBJ
         * @static
         * @param {message.IW_C_ADD_SIMOBJ} message W_C_ADD_SIMOBJ message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        W_C_ADD_SIMOBJ.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a W_C_ADD_SIMOBJ message from the specified reader or buffer.
         * @function decode
         * @memberof message.W_C_ADD_SIMOBJ
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {message.W_C_ADD_SIMOBJ} W_C_ADD_SIMOBJ
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        W_C_ADD_SIMOBJ.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.W_C_ADD_SIMOBJ();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.PacketHead = $root.message.Ipacket.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.Id = reader.int64();
                    break;
                case 3:
                    message.Pos = $root.message.Point3F.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.Rotation = reader.float();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a W_C_ADD_SIMOBJ message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof message.W_C_ADD_SIMOBJ
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {message.W_C_ADD_SIMOBJ} W_C_ADD_SIMOBJ
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        W_C_ADD_SIMOBJ.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a W_C_ADD_SIMOBJ message.
         * @function verify
         * @memberof message.W_C_ADD_SIMOBJ
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        W_C_ADD_SIMOBJ.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead")) {
                var error = $root.message.Ipacket.verify(message.PacketHead);
                if (error)
                    return "PacketHead." + error;
            }
            if (message.Id != null && message.hasOwnProperty("Id"))
                if (!$util.isInteger(message.Id) && !(message.Id && $util.isInteger(message.Id.low) && $util.isInteger(message.Id.high)))
                    return "Id: integer|Long expected";
            if (message.Pos != null && message.hasOwnProperty("Pos")) {
                var error = $root.message.Point3F.verify(message.Pos);
                if (error)
                    return "Pos." + error;
            }
            if (message.Rotation != null && message.hasOwnProperty("Rotation"))
                if (typeof message.Rotation !== "number")
                    return "Rotation: number expected";
            return null;
        };

        /**
         * Creates a W_C_ADD_SIMOBJ message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof message.W_C_ADD_SIMOBJ
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {message.W_C_ADD_SIMOBJ} W_C_ADD_SIMOBJ
         */
        W_C_ADD_SIMOBJ.fromObject = function fromObject(object) {
            if (object instanceof $root.message.W_C_ADD_SIMOBJ)
                return object;
            var message = new $root.message.W_C_ADD_SIMOBJ();
            if (object.PacketHead != null) {
                if (typeof object.PacketHead !== "object")
                    throw TypeError(".message.W_C_ADD_SIMOBJ.PacketHead: object expected");
                message.PacketHead = $root.message.Ipacket.fromObject(object.PacketHead);
            }
            if (object.Id != null)
                if ($util.Long)
                    (message.Id = $util.Long.fromValue(object.Id)).unsigned = false;
                else if (typeof object.Id === "string")
                    message.Id = parseInt(object.Id, 10);
                else if (typeof object.Id === "number")
                    message.Id = object.Id;
                else if (typeof object.Id === "object")
                    message.Id = new $util.LongBits(object.Id.low >>> 0, object.Id.high >>> 0).toNumber();
            if (object.Pos != null) {
                if (typeof object.Pos !== "object")
                    throw TypeError(".message.W_C_ADD_SIMOBJ.Pos: object expected");
                message.Pos = $root.message.Point3F.fromObject(object.Pos);
            }
            if (object.Rotation != null)
                message.Rotation = Number(object.Rotation);
            return message;
        };

        /**
         * Creates a plain object from a W_C_ADD_SIMOBJ message. Also converts values to other types if specified.
         * @function toObject
         * @memberof message.W_C_ADD_SIMOBJ
         * @static
         * @param {message.W_C_ADD_SIMOBJ} message W_C_ADD_SIMOBJ
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        W_C_ADD_SIMOBJ.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.defaults) {
                object.PacketHead = null;
                if ($util.Long) {
                    var long = new $util.Long(0, 0, false);
                    object.Id = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.Id = options.longs === String ? "0" : 0;
                object.Pos = null;
                object.Rotation = 0;
            }
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                object.PacketHead = $root.message.Ipacket.toObject(message.PacketHead, options);
            if (message.Id != null && message.hasOwnProperty("Id"))
                if (typeof message.Id === "number")
                    object.Id = options.longs === String ? String(message.Id) : message.Id;
                else
                    object.Id = options.longs === String ? $util.Long.prototype.toString.call(message.Id) : options.longs === Number ? new $util.LongBits(message.Id.low >>> 0, message.Id.high >>> 0).toNumber() : message.Id;
            if (message.Pos != null && message.hasOwnProperty("Pos"))
                object.Pos = $root.message.Point3F.toObject(message.Pos, options);
            if (message.Rotation != null && message.hasOwnProperty("Rotation"))
                object.Rotation = options.json && !isFinite(message.Rotation) ? String(message.Rotation) : message.Rotation;
            return object;
        };

        /**
         * Converts this W_C_ADD_SIMOBJ to JSON.
         * @function toJSON
         * @memberof message.W_C_ADD_SIMOBJ
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        W_C_ADD_SIMOBJ.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return W_C_ADD_SIMOBJ;
    })();

    message.C_W_LoginCopyMap = (function() {

        /**
         * Properties of a C_W_LoginCopyMap.
         * @memberof message
         * @interface IC_W_LoginCopyMap
         * @property {message.IIpacket|null} [PacketHead] C_W_LoginCopyMap PacketHead
         * @property {number|null} [DataId] C_W_LoginCopyMap DataId
         */

        /**
         * Constructs a new C_W_LoginCopyMap.
         * @memberof message
         * @classdesc Represents a C_W_LoginCopyMap.
         * @implements IC_W_LoginCopyMap
         * @constructor
         * @param {message.IC_W_LoginCopyMap=} [properties] Properties to set
         */
        function C_W_LoginCopyMap(properties) {
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * C_W_LoginCopyMap PacketHead.
         * @member {message.IIpacket|null|undefined} PacketHead
         * @memberof message.C_W_LoginCopyMap
         * @instance
         */
        C_W_LoginCopyMap.prototype.PacketHead = null;

        /**
         * C_W_LoginCopyMap DataId.
         * @member {number} DataId
         * @memberof message.C_W_LoginCopyMap
         * @instance
         */
        C_W_LoginCopyMap.prototype.DataId = 0;

        /**
         * Creates a new C_W_LoginCopyMap instance using the specified properties.
         * @function create
         * @memberof message.C_W_LoginCopyMap
         * @static
         * @param {message.IC_W_LoginCopyMap=} [properties] Properties to set
         * @returns {message.C_W_LoginCopyMap} C_W_LoginCopyMap instance
         */
        C_W_LoginCopyMap.create = function create(properties) {
            return new C_W_LoginCopyMap(properties);
        };

        /**
         * Encodes the specified C_W_LoginCopyMap message. Does not implicitly {@link message.C_W_LoginCopyMap.verify|verify} messages.
         * @function encode
         * @memberof message.C_W_LoginCopyMap
         * @static
         * @param {message.IC_W_LoginCopyMap} message C_W_LoginCopyMap message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        C_W_LoginCopyMap.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                $root.message.Ipacket.encode(message.PacketHead, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
            if (message.DataId != null && message.hasOwnProperty("DataId"))
                writer.uint32(/* id 2, wireType 0 =*/16).int32(message.DataId);
            return writer;
        };

        /**
         * Encodes the specified C_W_LoginCopyMap message, length delimited. Does not implicitly {@link message.C_W_LoginCopyMap.verify|verify} messages.
         * @function encodeDelimited
         * @memberof message.C_W_LoginCopyMap
         * @static
         * @param {message.IC_W_LoginCopyMap} message C_W_LoginCopyMap message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        C_W_LoginCopyMap.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a C_W_LoginCopyMap message from the specified reader or buffer.
         * @function decode
         * @memberof message.C_W_LoginCopyMap
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {message.C_W_LoginCopyMap} C_W_LoginCopyMap
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        C_W_LoginCopyMap.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.C_W_LoginCopyMap();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.PacketHead = $root.message.Ipacket.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.DataId = reader.int32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a C_W_LoginCopyMap message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof message.C_W_LoginCopyMap
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {message.C_W_LoginCopyMap} C_W_LoginCopyMap
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        C_W_LoginCopyMap.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a C_W_LoginCopyMap message.
         * @function verify
         * @memberof message.C_W_LoginCopyMap
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        C_W_LoginCopyMap.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead")) {
                var error = $root.message.Ipacket.verify(message.PacketHead);
                if (error)
                    return "PacketHead." + error;
            }
            if (message.DataId != null && message.hasOwnProperty("DataId"))
                if (!$util.isInteger(message.DataId))
                    return "DataId: integer expected";
            return null;
        };

        /**
         * Creates a C_W_LoginCopyMap message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof message.C_W_LoginCopyMap
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {message.C_W_LoginCopyMap} C_W_LoginCopyMap
         */
        C_W_LoginCopyMap.fromObject = function fromObject(object) {
            if (object instanceof $root.message.C_W_LoginCopyMap)
                return object;
            var message = new $root.message.C_W_LoginCopyMap();
            if (object.PacketHead != null) {
                if (typeof object.PacketHead !== "object")
                    throw TypeError(".message.C_W_LoginCopyMap.PacketHead: object expected");
                message.PacketHead = $root.message.Ipacket.fromObject(object.PacketHead);
            }
            if (object.DataId != null)
                message.DataId = object.DataId | 0;
            return message;
        };

        /**
         * Creates a plain object from a C_W_LoginCopyMap message. Also converts values to other types if specified.
         * @function toObject
         * @memberof message.C_W_LoginCopyMap
         * @static
         * @param {message.C_W_LoginCopyMap} message C_W_LoginCopyMap
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        C_W_LoginCopyMap.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.defaults) {
                object.PacketHead = null;
                object.DataId = 0;
            }
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                object.PacketHead = $root.message.Ipacket.toObject(message.PacketHead, options);
            if (message.DataId != null && message.hasOwnProperty("DataId"))
                object.DataId = message.DataId;
            return object;
        };

        /**
         * Converts this C_W_LoginCopyMap to JSON.
         * @function toJSON
         * @memberof message.C_W_LoginCopyMap
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        C_W_LoginCopyMap.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return C_W_LoginCopyMap;
    })();

    /**
     * SERVICE enum.
     * @name message.SERVICE
     * @enum {string}
     * @property {number} NONE=0 NONE value
     * @property {number} CLIENT=1 CLIENT value
     * @property {number} GATESERVER=2 GATESERVER value
     * @property {number} ACCOUNTSERVER=3 ACCOUNTSERVER value
     * @property {number} WORLDSERVER=4 WORLDSERVER value
     * @property {number} ZONESERVER=5 ZONESERVER value
     * @property {number} WORLDDBSERVER=6 WORLDDBSERVER value
     */
    message.SERVICE = (function() {
        var valuesById = {}, values = Object.create(valuesById);
        values[valuesById[0] = "NONE"] = 0;
        values[valuesById[1] = "CLIENT"] = 1;
        values[valuesById[2] = "GATESERVER"] = 2;
        values[valuesById[3] = "ACCOUNTSERVER"] = 3;
        values[valuesById[4] = "WORLDSERVER"] = 4;
        values[valuesById[5] = "ZONESERVER"] = 5;
        values[valuesById[6] = "WORLDDBSERVER"] = 6;
        return values;
    })();

    /**
     * CHAT enum.
     * @name message.CHAT
     * @enum {string}
     * @property {number} MSG_TYPE_WORLD=0 MSG_TYPE_WORLD value
     * @property {number} MSG_TYPE_PRIVATE=1 MSG_TYPE_PRIVATE value
     * @property {number} MSG_TYPE_ORG=2 MSG_TYPE_ORG value
     * @property {number} MSG_TYPE_COUNT=3 MSG_TYPE_COUNT value
     */
    message.CHAT = (function() {
        var valuesById = {}, values = Object.create(valuesById);
        values[valuesById[0] = "MSG_TYPE_WORLD"] = 0;
        values[valuesById[1] = "MSG_TYPE_PRIVATE"] = 1;
        values[valuesById[2] = "MSG_TYPE_ORG"] = 2;
        values[valuesById[3] = "MSG_TYPE_COUNT"] = 3;
        return values;
    })();

    message.Ipacket = (function() {

        /**
         * Properties of an Ipacket.
         * @memberof message
         * @interface IIpacket
         * @property {number|null} [Stx] Ipacket Stx
         * @property {number|null} [DestServerType] Ipacket DestServerType
         * @property {number|null} [Ckx] Ipacket Ckx
         * @property {number|Long|null} [Id] Ipacket Id
         */

        /**
         * Constructs a new Ipacket.
         * @memberof message
         * @classdesc Represents an Ipacket.
         * @implements IIpacket
         * @constructor
         * @param {message.IIpacket=} [properties] Properties to set
         */
        function Ipacket(properties) {
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * Ipacket Stx.
         * @member {number} Stx
         * @memberof message.Ipacket
         * @instance
         */
        Ipacket.prototype.Stx = 0;

        /**
         * Ipacket DestServerType.
         * @member {number} DestServerType
         * @memberof message.Ipacket
         * @instance
         */
        Ipacket.prototype.DestServerType = 0;

        /**
         * Ipacket Ckx.
         * @member {number} Ckx
         * @memberof message.Ipacket
         * @instance
         */
        Ipacket.prototype.Ckx = 0;

        /**
         * Ipacket Id.
         * @member {number|Long} Id
         * @memberof message.Ipacket
         * @instance
         */
        Ipacket.prototype.Id = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * Creates a new Ipacket instance using the specified properties.
         * @function create
         * @memberof message.Ipacket
         * @static
         * @param {message.IIpacket=} [properties] Properties to set
         * @returns {message.Ipacket} Ipacket instance
         */
        Ipacket.create = function create(properties) {
            return new Ipacket(properties);
        };

        /**
         * Encodes the specified Ipacket message. Does not implicitly {@link message.Ipacket.verify|verify} messages.
         * @function encode
         * @memberof message.Ipacket
         * @static
         * @param {message.IIpacket} message Ipacket message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        Ipacket.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.Stx != null && message.hasOwnProperty("Stx"))
                writer.uint32(/* id 1, wireType 0 =*/8).int32(message.Stx);
            if (message.DestServerType != null && message.hasOwnProperty("DestServerType"))
                writer.uint32(/* id 2, wireType 0 =*/16).int32(message.DestServerType);
            if (message.Ckx != null && message.hasOwnProperty("Ckx"))
                writer.uint32(/* id 3, wireType 0 =*/24).int32(message.Ckx);
            if (message.Id != null && message.hasOwnProperty("Id"))
                writer.uint32(/* id 4, wireType 0 =*/32).int64(message.Id);
            return writer;
        };

        /**
         * Encodes the specified Ipacket message, length delimited. Does not implicitly {@link message.Ipacket.verify|verify} messages.
         * @function encodeDelimited
         * @memberof message.Ipacket
         * @static
         * @param {message.IIpacket} message Ipacket message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        Ipacket.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes an Ipacket message from the specified reader or buffer.
         * @function decode
         * @memberof message.Ipacket
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {message.Ipacket} Ipacket
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        Ipacket.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.Ipacket();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.Stx = reader.int32();
                    break;
                case 2:
                    message.DestServerType = reader.int32();
                    break;
                case 3:
                    message.Ckx = reader.int32();
                    break;
                case 4:
                    message.Id = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes an Ipacket message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof message.Ipacket
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {message.Ipacket} Ipacket
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        Ipacket.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies an Ipacket message.
         * @function verify
         * @memberof message.Ipacket
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        Ipacket.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.Stx != null && message.hasOwnProperty("Stx"))
                if (!$util.isInteger(message.Stx))
                    return "Stx: integer expected";
            if (message.DestServerType != null && message.hasOwnProperty("DestServerType"))
                if (!$util.isInteger(message.DestServerType))
                    return "DestServerType: integer expected";
            if (message.Ckx != null && message.hasOwnProperty("Ckx"))
                if (!$util.isInteger(message.Ckx))
                    return "Ckx: integer expected";
            if (message.Id != null && message.hasOwnProperty("Id"))
                if (!$util.isInteger(message.Id) && !(message.Id && $util.isInteger(message.Id.low) && $util.isInteger(message.Id.high)))
                    return "Id: integer|Long expected";
            return null;
        };

        /**
         * Creates an Ipacket message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof message.Ipacket
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {message.Ipacket} Ipacket
         */
        Ipacket.fromObject = function fromObject(object) {
            if (object instanceof $root.message.Ipacket)
                return object;
            var message = new $root.message.Ipacket();
            if (object.Stx != null)
                message.Stx = object.Stx | 0;
            if (object.DestServerType != null)
                message.DestServerType = object.DestServerType | 0;
            if (object.Ckx != null)
                message.Ckx = object.Ckx | 0;
            if (object.Id != null)
                if ($util.Long)
                    (message.Id = $util.Long.fromValue(object.Id)).unsigned = false;
                else if (typeof object.Id === "string")
                    message.Id = parseInt(object.Id, 10);
                else if (typeof object.Id === "number")
                    message.Id = object.Id;
                else if (typeof object.Id === "object")
                    message.Id = new $util.LongBits(object.Id.low >>> 0, object.Id.high >>> 0).toNumber();
            return message;
        };

        /**
         * Creates a plain object from an Ipacket message. Also converts values to other types if specified.
         * @function toObject
         * @memberof message.Ipacket
         * @static
         * @param {message.Ipacket} message Ipacket
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        Ipacket.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.defaults) {
                object.Stx = 0;
                object.DestServerType = 0;
                object.Ckx = 0;
                if ($util.Long) {
                    var long = new $util.Long(0, 0, false);
                    object.Id = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.Id = options.longs === String ? "0" : 0;
            }
            if (message.Stx != null && message.hasOwnProperty("Stx"))
                object.Stx = message.Stx;
            if (message.DestServerType != null && message.hasOwnProperty("DestServerType"))
                object.DestServerType = message.DestServerType;
            if (message.Ckx != null && message.hasOwnProperty("Ckx"))
                object.Ckx = message.Ckx;
            if (message.Id != null && message.hasOwnProperty("Id"))
                if (typeof message.Id === "number")
                    object.Id = options.longs === String ? String(message.Id) : message.Id;
                else
                    object.Id = options.longs === String ? $util.Long.prototype.toString.call(message.Id) : options.longs === Number ? new $util.LongBits(message.Id.low >>> 0, message.Id.high >>> 0).toNumber() : message.Id;
            return object;
        };

        /**
         * Converts this Ipacket to JSON.
         * @function toJSON
         * @memberof message.Ipacket
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        Ipacket.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return Ipacket;
    })();

    message.PlayerData = (function() {

        /**
         * Properties of a PlayerData.
         * @memberof message
         * @interface IPlayerData
         * @property {number|Long|null} [PlayerID] PlayerData PlayerID
         * @property {string|null} [PlayerName] PlayerData PlayerName
         * @property {number|null} [PlayerGold] PlayerData PlayerGold
         */

        /**
         * Constructs a new PlayerData.
         * @memberof message
         * @classdesc Represents a PlayerData.
         * @implements IPlayerData
         * @constructor
         * @param {message.IPlayerData=} [properties] Properties to set
         */
        function PlayerData(properties) {
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * PlayerData PlayerID.
         * @member {number|Long} PlayerID
         * @memberof message.PlayerData
         * @instance
         */
        PlayerData.prototype.PlayerID = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * PlayerData PlayerName.
         * @member {string} PlayerName
         * @memberof message.PlayerData
         * @instance
         */
        PlayerData.prototype.PlayerName = "";

        /**
         * PlayerData PlayerGold.
         * @member {number} PlayerGold
         * @memberof message.PlayerData
         * @instance
         */
        PlayerData.prototype.PlayerGold = 0;

        /**
         * Creates a new PlayerData instance using the specified properties.
         * @function create
         * @memberof message.PlayerData
         * @static
         * @param {message.IPlayerData=} [properties] Properties to set
         * @returns {message.PlayerData} PlayerData instance
         */
        PlayerData.create = function create(properties) {
            return new PlayerData(properties);
        };

        /**
         * Encodes the specified PlayerData message. Does not implicitly {@link message.PlayerData.verify|verify} messages.
         * @function encode
         * @memberof message.PlayerData
         * @static
         * @param {message.IPlayerData} message PlayerData message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        PlayerData.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.PlayerID != null && message.hasOwnProperty("PlayerID"))
                writer.uint32(/* id 1, wireType 0 =*/8).int64(message.PlayerID);
            if (message.PlayerName != null && message.hasOwnProperty("PlayerName"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.PlayerName);
            if (message.PlayerGold != null && message.hasOwnProperty("PlayerGold"))
                writer.uint32(/* id 3, wireType 0 =*/24).int32(message.PlayerGold);
            return writer;
        };

        /**
         * Encodes the specified PlayerData message, length delimited. Does not implicitly {@link message.PlayerData.verify|verify} messages.
         * @function encodeDelimited
         * @memberof message.PlayerData
         * @static
         * @param {message.IPlayerData} message PlayerData message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        PlayerData.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a PlayerData message from the specified reader or buffer.
         * @function decode
         * @memberof message.PlayerData
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {message.PlayerData} PlayerData
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        PlayerData.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.PlayerData();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.PlayerID = reader.int64();
                    break;
                case 2:
                    message.PlayerName = reader.string();
                    break;
                case 3:
                    message.PlayerGold = reader.int32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a PlayerData message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof message.PlayerData
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {message.PlayerData} PlayerData
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        PlayerData.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a PlayerData message.
         * @function verify
         * @memberof message.PlayerData
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        PlayerData.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.PlayerID != null && message.hasOwnProperty("PlayerID"))
                if (!$util.isInteger(message.PlayerID) && !(message.PlayerID && $util.isInteger(message.PlayerID.low) && $util.isInteger(message.PlayerID.high)))
                    return "PlayerID: integer|Long expected";
            if (message.PlayerName != null && message.hasOwnProperty("PlayerName"))
                if (!$util.isString(message.PlayerName))
                    return "PlayerName: string expected";
            if (message.PlayerGold != null && message.hasOwnProperty("PlayerGold"))
                if (!$util.isInteger(message.PlayerGold))
                    return "PlayerGold: integer expected";
            return null;
        };

        /**
         * Creates a PlayerData message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof message.PlayerData
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {message.PlayerData} PlayerData
         */
        PlayerData.fromObject = function fromObject(object) {
            if (object instanceof $root.message.PlayerData)
                return object;
            var message = new $root.message.PlayerData();
            if (object.PlayerID != null)
                if ($util.Long)
                    (message.PlayerID = $util.Long.fromValue(object.PlayerID)).unsigned = false;
                else if (typeof object.PlayerID === "string")
                    message.PlayerID = parseInt(object.PlayerID, 10);
                else if (typeof object.PlayerID === "number")
                    message.PlayerID = object.PlayerID;
                else if (typeof object.PlayerID === "object")
                    message.PlayerID = new $util.LongBits(object.PlayerID.low >>> 0, object.PlayerID.high >>> 0).toNumber();
            if (object.PlayerName != null)
                message.PlayerName = String(object.PlayerName);
            if (object.PlayerGold != null)
                message.PlayerGold = object.PlayerGold | 0;
            return message;
        };

        /**
         * Creates a plain object from a PlayerData message. Also converts values to other types if specified.
         * @function toObject
         * @memberof message.PlayerData
         * @static
         * @param {message.PlayerData} message PlayerData
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        PlayerData.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.defaults) {
                if ($util.Long) {
                    var long = new $util.Long(0, 0, false);
                    object.PlayerID = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.PlayerID = options.longs === String ? "0" : 0;
                object.PlayerName = "";
                object.PlayerGold = 0;
            }
            if (message.PlayerID != null && message.hasOwnProperty("PlayerID"))
                if (typeof message.PlayerID === "number")
                    object.PlayerID = options.longs === String ? String(message.PlayerID) : message.PlayerID;
                else
                    object.PlayerID = options.longs === String ? $util.Long.prototype.toString.call(message.PlayerID) : options.longs === Number ? new $util.LongBits(message.PlayerID.low >>> 0, message.PlayerID.high >>> 0).toNumber() : message.PlayerID;
            if (message.PlayerName != null && message.hasOwnProperty("PlayerName"))
                object.PlayerName = message.PlayerName;
            if (message.PlayerGold != null && message.hasOwnProperty("PlayerGold"))
                object.PlayerGold = message.PlayerGold;
            return object;
        };

        /**
         * Converts this PlayerData to JSON.
         * @function toJSON
         * @memberof message.PlayerData
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        PlayerData.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return PlayerData;
    })();

    message.W_C_CreatePlayerResponse = (function() {

        /**
         * Properties of a W_C_CreatePlayerResponse.
         * @memberof message
         * @interface IW_C_CreatePlayerResponse
         * @property {message.IIpacket|null} [PacketHead] W_C_CreatePlayerResponse PacketHead
         * @property {number|null} [Error] W_C_CreatePlayerResponse Error
         * @property {number|Long|null} [PlayerId] W_C_CreatePlayerResponse PlayerId
         */

        /**
         * Constructs a new W_C_CreatePlayerResponse.
         * @memberof message
         * @classdesc Represents a W_C_CreatePlayerResponse.
         * @implements IW_C_CreatePlayerResponse
         * @constructor
         * @param {message.IW_C_CreatePlayerResponse=} [properties] Properties to set
         */
        function W_C_CreatePlayerResponse(properties) {
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * W_C_CreatePlayerResponse PacketHead.
         * @member {message.IIpacket|null|undefined} PacketHead
         * @memberof message.W_C_CreatePlayerResponse
         * @instance
         */
        W_C_CreatePlayerResponse.prototype.PacketHead = null;

        /**
         * W_C_CreatePlayerResponse Error.
         * @member {number} Error
         * @memberof message.W_C_CreatePlayerResponse
         * @instance
         */
        W_C_CreatePlayerResponse.prototype.Error = 0;

        /**
         * W_C_CreatePlayerResponse PlayerId.
         * @member {number|Long} PlayerId
         * @memberof message.W_C_CreatePlayerResponse
         * @instance
         */
        W_C_CreatePlayerResponse.prototype.PlayerId = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * Creates a new W_C_CreatePlayerResponse instance using the specified properties.
         * @function create
         * @memberof message.W_C_CreatePlayerResponse
         * @static
         * @param {message.IW_C_CreatePlayerResponse=} [properties] Properties to set
         * @returns {message.W_C_CreatePlayerResponse} W_C_CreatePlayerResponse instance
         */
        W_C_CreatePlayerResponse.create = function create(properties) {
            return new W_C_CreatePlayerResponse(properties);
        };

        /**
         * Encodes the specified W_C_CreatePlayerResponse message. Does not implicitly {@link message.W_C_CreatePlayerResponse.verify|verify} messages.
         * @function encode
         * @memberof message.W_C_CreatePlayerResponse
         * @static
         * @param {message.IW_C_CreatePlayerResponse} message W_C_CreatePlayerResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        W_C_CreatePlayerResponse.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                $root.message.Ipacket.encode(message.PacketHead, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
            if (message.Error != null && message.hasOwnProperty("Error"))
                writer.uint32(/* id 2, wireType 0 =*/16).int32(message.Error);
            if (message.PlayerId != null && message.hasOwnProperty("PlayerId"))
                writer.uint32(/* id 3, wireType 0 =*/24).int64(message.PlayerId);
            return writer;
        };

        /**
         * Encodes the specified W_C_CreatePlayerResponse message, length delimited. Does not implicitly {@link message.W_C_CreatePlayerResponse.verify|verify} messages.
         * @function encodeDelimited
         * @memberof message.W_C_CreatePlayerResponse
         * @static
         * @param {message.IW_C_CreatePlayerResponse} message W_C_CreatePlayerResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        W_C_CreatePlayerResponse.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a W_C_CreatePlayerResponse message from the specified reader or buffer.
         * @function decode
         * @memberof message.W_C_CreatePlayerResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {message.W_C_CreatePlayerResponse} W_C_CreatePlayerResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        W_C_CreatePlayerResponse.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.W_C_CreatePlayerResponse();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.PacketHead = $root.message.Ipacket.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.Error = reader.int32();
                    break;
                case 3:
                    message.PlayerId = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a W_C_CreatePlayerResponse message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof message.W_C_CreatePlayerResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {message.W_C_CreatePlayerResponse} W_C_CreatePlayerResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        W_C_CreatePlayerResponse.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a W_C_CreatePlayerResponse message.
         * @function verify
         * @memberof message.W_C_CreatePlayerResponse
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        W_C_CreatePlayerResponse.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead")) {
                var error = $root.message.Ipacket.verify(message.PacketHead);
                if (error)
                    return "PacketHead." + error;
            }
            if (message.Error != null && message.hasOwnProperty("Error"))
                if (!$util.isInteger(message.Error))
                    return "Error: integer expected";
            if (message.PlayerId != null && message.hasOwnProperty("PlayerId"))
                if (!$util.isInteger(message.PlayerId) && !(message.PlayerId && $util.isInteger(message.PlayerId.low) && $util.isInteger(message.PlayerId.high)))
                    return "PlayerId: integer|Long expected";
            return null;
        };

        /**
         * Creates a W_C_CreatePlayerResponse message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof message.W_C_CreatePlayerResponse
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {message.W_C_CreatePlayerResponse} W_C_CreatePlayerResponse
         */
        W_C_CreatePlayerResponse.fromObject = function fromObject(object) {
            if (object instanceof $root.message.W_C_CreatePlayerResponse)
                return object;
            var message = new $root.message.W_C_CreatePlayerResponse();
            if (object.PacketHead != null) {
                if (typeof object.PacketHead !== "object")
                    throw TypeError(".message.W_C_CreatePlayerResponse.PacketHead: object expected");
                message.PacketHead = $root.message.Ipacket.fromObject(object.PacketHead);
            }
            if (object.Error != null)
                message.Error = object.Error | 0;
            if (object.PlayerId != null)
                if ($util.Long)
                    (message.PlayerId = $util.Long.fromValue(object.PlayerId)).unsigned = false;
                else if (typeof object.PlayerId === "string")
                    message.PlayerId = parseInt(object.PlayerId, 10);
                else if (typeof object.PlayerId === "number")
                    message.PlayerId = object.PlayerId;
                else if (typeof object.PlayerId === "object")
                    message.PlayerId = new $util.LongBits(object.PlayerId.low >>> 0, object.PlayerId.high >>> 0).toNumber();
            return message;
        };

        /**
         * Creates a plain object from a W_C_CreatePlayerResponse message. Also converts values to other types if specified.
         * @function toObject
         * @memberof message.W_C_CreatePlayerResponse
         * @static
         * @param {message.W_C_CreatePlayerResponse} message W_C_CreatePlayerResponse
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        W_C_CreatePlayerResponse.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.defaults) {
                object.PacketHead = null;
                object.Error = 0;
                if ($util.Long) {
                    var long = new $util.Long(0, 0, false);
                    object.PlayerId = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.PlayerId = options.longs === String ? "0" : 0;
            }
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                object.PacketHead = $root.message.Ipacket.toObject(message.PacketHead, options);
            if (message.Error != null && message.hasOwnProperty("Error"))
                object.Error = message.Error;
            if (message.PlayerId != null && message.hasOwnProperty("PlayerId"))
                if (typeof message.PlayerId === "number")
                    object.PlayerId = options.longs === String ? String(message.PlayerId) : message.PlayerId;
                else
                    object.PlayerId = options.longs === String ? $util.Long.prototype.toString.call(message.PlayerId) : options.longs === Number ? new $util.LongBits(message.PlayerId.low >>> 0, message.PlayerId.high >>> 0).toNumber() : message.PlayerId;
            return object;
        };

        /**
         * Converts this W_C_CreatePlayerResponse to JSON.
         * @function toJSON
         * @memberof message.W_C_CreatePlayerResponse
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        W_C_CreatePlayerResponse.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return W_C_CreatePlayerResponse;
    })();

    message.W_C_SelectPlayerResponse = (function() {

        /**
         * Properties of a W_C_SelectPlayerResponse.
         * @memberof message
         * @interface IW_C_SelectPlayerResponse
         * @property {message.IIpacket|null} [PacketHead] W_C_SelectPlayerResponse PacketHead
         * @property {number|Long|null} [AccountId] W_C_SelectPlayerResponse AccountId
         * @property {Array.<message.IPlayerData>|null} [PlayerData] W_C_SelectPlayerResponse PlayerData
         */

        /**
         * Constructs a new W_C_SelectPlayerResponse.
         * @memberof message
         * @classdesc Represents a W_C_SelectPlayerResponse.
         * @implements IW_C_SelectPlayerResponse
         * @constructor
         * @param {message.IW_C_SelectPlayerResponse=} [properties] Properties to set
         */
        function W_C_SelectPlayerResponse(properties) {
            this.PlayerData = [];
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * W_C_SelectPlayerResponse PacketHead.
         * @member {message.IIpacket|null|undefined} PacketHead
         * @memberof message.W_C_SelectPlayerResponse
         * @instance
         */
        W_C_SelectPlayerResponse.prototype.PacketHead = null;

        /**
         * W_C_SelectPlayerResponse AccountId.
         * @member {number|Long} AccountId
         * @memberof message.W_C_SelectPlayerResponse
         * @instance
         */
        W_C_SelectPlayerResponse.prototype.AccountId = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * W_C_SelectPlayerResponse PlayerData.
         * @member {Array.<message.IPlayerData>} PlayerData
         * @memberof message.W_C_SelectPlayerResponse
         * @instance
         */
        W_C_SelectPlayerResponse.prototype.PlayerData = $util.emptyArray;

        /**
         * Creates a new W_C_SelectPlayerResponse instance using the specified properties.
         * @function create
         * @memberof message.W_C_SelectPlayerResponse
         * @static
         * @param {message.IW_C_SelectPlayerResponse=} [properties] Properties to set
         * @returns {message.W_C_SelectPlayerResponse} W_C_SelectPlayerResponse instance
         */
        W_C_SelectPlayerResponse.create = function create(properties) {
            return new W_C_SelectPlayerResponse(properties);
        };

        /**
         * Encodes the specified W_C_SelectPlayerResponse message. Does not implicitly {@link message.W_C_SelectPlayerResponse.verify|verify} messages.
         * @function encode
         * @memberof message.W_C_SelectPlayerResponse
         * @static
         * @param {message.IW_C_SelectPlayerResponse} message W_C_SelectPlayerResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        W_C_SelectPlayerResponse.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                $root.message.Ipacket.encode(message.PacketHead, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
            if (message.AccountId != null && message.hasOwnProperty("AccountId"))
                writer.uint32(/* id 2, wireType 0 =*/16).int64(message.AccountId);
            if (message.PlayerData != null && message.PlayerData.length)
                for (var i = 0; i < message.PlayerData.length; ++i)
                    $root.message.PlayerData.encode(message.PlayerData[i], writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
            return writer;
        };

        /**
         * Encodes the specified W_C_SelectPlayerResponse message, length delimited. Does not implicitly {@link message.W_C_SelectPlayerResponse.verify|verify} messages.
         * @function encodeDelimited
         * @memberof message.W_C_SelectPlayerResponse
         * @static
         * @param {message.IW_C_SelectPlayerResponse} message W_C_SelectPlayerResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        W_C_SelectPlayerResponse.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a W_C_SelectPlayerResponse message from the specified reader or buffer.
         * @function decode
         * @memberof message.W_C_SelectPlayerResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {message.W_C_SelectPlayerResponse} W_C_SelectPlayerResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        W_C_SelectPlayerResponse.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.W_C_SelectPlayerResponse();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.PacketHead = $root.message.Ipacket.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.AccountId = reader.int64();
                    break;
                case 3:
                    if (!(message.PlayerData && message.PlayerData.length))
                        message.PlayerData = [];
                    message.PlayerData.push($root.message.PlayerData.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a W_C_SelectPlayerResponse message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof message.W_C_SelectPlayerResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {message.W_C_SelectPlayerResponse} W_C_SelectPlayerResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        W_C_SelectPlayerResponse.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a W_C_SelectPlayerResponse message.
         * @function verify
         * @memberof message.W_C_SelectPlayerResponse
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        W_C_SelectPlayerResponse.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead")) {
                var error = $root.message.Ipacket.verify(message.PacketHead);
                if (error)
                    return "PacketHead." + error;
            }
            if (message.AccountId != null && message.hasOwnProperty("AccountId"))
                if (!$util.isInteger(message.AccountId) && !(message.AccountId && $util.isInteger(message.AccountId.low) && $util.isInteger(message.AccountId.high)))
                    return "AccountId: integer|Long expected";
            if (message.PlayerData != null && message.hasOwnProperty("PlayerData")) {
                if (!Array.isArray(message.PlayerData))
                    return "PlayerData: array expected";
                for (var i = 0; i < message.PlayerData.length; ++i) {
                    var error = $root.message.PlayerData.verify(message.PlayerData[i]);
                    if (error)
                        return "PlayerData." + error;
                }
            }
            return null;
        };

        /**
         * Creates a W_C_SelectPlayerResponse message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof message.W_C_SelectPlayerResponse
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {message.W_C_SelectPlayerResponse} W_C_SelectPlayerResponse
         */
        W_C_SelectPlayerResponse.fromObject = function fromObject(object) {
            if (object instanceof $root.message.W_C_SelectPlayerResponse)
                return object;
            var message = new $root.message.W_C_SelectPlayerResponse();
            if (object.PacketHead != null) {
                if (typeof object.PacketHead !== "object")
                    throw TypeError(".message.W_C_SelectPlayerResponse.PacketHead: object expected");
                message.PacketHead = $root.message.Ipacket.fromObject(object.PacketHead);
            }
            if (object.AccountId != null)
                if ($util.Long)
                    (message.AccountId = $util.Long.fromValue(object.AccountId)).unsigned = false;
                else if (typeof object.AccountId === "string")
                    message.AccountId = parseInt(object.AccountId, 10);
                else if (typeof object.AccountId === "number")
                    message.AccountId = object.AccountId;
                else if (typeof object.AccountId === "object")
                    message.AccountId = new $util.LongBits(object.AccountId.low >>> 0, object.AccountId.high >>> 0).toNumber();
            if (object.PlayerData) {
                if (!Array.isArray(object.PlayerData))
                    throw TypeError(".message.W_C_SelectPlayerResponse.PlayerData: array expected");
                message.PlayerData = [];
                for (var i = 0; i < object.PlayerData.length; ++i) {
                    if (typeof object.PlayerData[i] !== "object")
                        throw TypeError(".message.W_C_SelectPlayerResponse.PlayerData: object expected");
                    message.PlayerData[i] = $root.message.PlayerData.fromObject(object.PlayerData[i]);
                }
            }
            return message;
        };

        /**
         * Creates a plain object from a W_C_SelectPlayerResponse message. Also converts values to other types if specified.
         * @function toObject
         * @memberof message.W_C_SelectPlayerResponse
         * @static
         * @param {message.W_C_SelectPlayerResponse} message W_C_SelectPlayerResponse
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        W_C_SelectPlayerResponse.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.arrays || options.defaults)
                object.PlayerData = [];
            if (options.defaults) {
                object.PacketHead = null;
                if ($util.Long) {
                    var long = new $util.Long(0, 0, false);
                    object.AccountId = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.AccountId = options.longs === String ? "0" : 0;
            }
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                object.PacketHead = $root.message.Ipacket.toObject(message.PacketHead, options);
            if (message.AccountId != null && message.hasOwnProperty("AccountId"))
                if (typeof message.AccountId === "number")
                    object.AccountId = options.longs === String ? String(message.AccountId) : message.AccountId;
                else
                    object.AccountId = options.longs === String ? $util.Long.prototype.toString.call(message.AccountId) : options.longs === Number ? new $util.LongBits(message.AccountId.low >>> 0, message.AccountId.high >>> 0).toNumber() : message.AccountId;
            if (message.PlayerData && message.PlayerData.length) {
                object.PlayerData = [];
                for (var j = 0; j < message.PlayerData.length; ++j)
                    object.PlayerData[j] = $root.message.PlayerData.toObject(message.PlayerData[j], options);
            }
            return object;
        };

        /**
         * Converts this W_C_SelectPlayerResponse to JSON.
         * @function toJSON
         * @memberof message.W_C_SelectPlayerResponse
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        W_C_SelectPlayerResponse.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return W_C_SelectPlayerResponse;
    })();

    message.A_C_RegisterResponse = (function() {

        /**
         * Properties of a A_C_RegisterResponse.
         * @memberof message
         * @interface IA_C_RegisterResponse
         * @property {message.IIpacket|null} [PacketHead] A_C_RegisterResponse PacketHead
         * @property {number|null} [Error] A_C_RegisterResponse Error
         * @property {number|null} [SocketId] A_C_RegisterResponse SocketId
         */

        /**
         * Constructs a new A_C_RegisterResponse.
         * @memberof message
         * @classdesc Represents a A_C_RegisterResponse.
         * @implements IA_C_RegisterResponse
         * @constructor
         * @param {message.IA_C_RegisterResponse=} [properties] Properties to set
         */
        function A_C_RegisterResponse(properties) {
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * A_C_RegisterResponse PacketHead.
         * @member {message.IIpacket|null|undefined} PacketHead
         * @memberof message.A_C_RegisterResponse
         * @instance
         */
        A_C_RegisterResponse.prototype.PacketHead = null;

        /**
         * A_C_RegisterResponse Error.
         * @member {number} Error
         * @memberof message.A_C_RegisterResponse
         * @instance
         */
        A_C_RegisterResponse.prototype.Error = 0;

        /**
         * A_C_RegisterResponse SocketId.
         * @member {number} SocketId
         * @memberof message.A_C_RegisterResponse
         * @instance
         */
        A_C_RegisterResponse.prototype.SocketId = 0;

        /**
         * Creates a new A_C_RegisterResponse instance using the specified properties.
         * @function create
         * @memberof message.A_C_RegisterResponse
         * @static
         * @param {message.IA_C_RegisterResponse=} [properties] Properties to set
         * @returns {message.A_C_RegisterResponse} A_C_RegisterResponse instance
         */
        A_C_RegisterResponse.create = function create(properties) {
            return new A_C_RegisterResponse(properties);
        };

        /**
         * Encodes the specified A_C_RegisterResponse message. Does not implicitly {@link message.A_C_RegisterResponse.verify|verify} messages.
         * @function encode
         * @memberof message.A_C_RegisterResponse
         * @static
         * @param {message.IA_C_RegisterResponse} message A_C_RegisterResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        A_C_RegisterResponse.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                $root.message.Ipacket.encode(message.PacketHead, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
            if (message.Error != null && message.hasOwnProperty("Error"))
                writer.uint32(/* id 2, wireType 0 =*/16).int32(message.Error);
            if (message.SocketId != null && message.hasOwnProperty("SocketId"))
                writer.uint32(/* id 3, wireType 0 =*/24).int32(message.SocketId);
            return writer;
        };

        /**
         * Encodes the specified A_C_RegisterResponse message, length delimited. Does not implicitly {@link message.A_C_RegisterResponse.verify|verify} messages.
         * @function encodeDelimited
         * @memberof message.A_C_RegisterResponse
         * @static
         * @param {message.IA_C_RegisterResponse} message A_C_RegisterResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        A_C_RegisterResponse.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a A_C_RegisterResponse message from the specified reader or buffer.
         * @function decode
         * @memberof message.A_C_RegisterResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {message.A_C_RegisterResponse} A_C_RegisterResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        A_C_RegisterResponse.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.A_C_RegisterResponse();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.PacketHead = $root.message.Ipacket.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.Error = reader.int32();
                    break;
                case 3:
                    message.SocketId = reader.int32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a A_C_RegisterResponse message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof message.A_C_RegisterResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {message.A_C_RegisterResponse} A_C_RegisterResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        A_C_RegisterResponse.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a A_C_RegisterResponse message.
         * @function verify
         * @memberof message.A_C_RegisterResponse
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        A_C_RegisterResponse.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead")) {
                var error = $root.message.Ipacket.verify(message.PacketHead);
                if (error)
                    return "PacketHead." + error;
            }
            if (message.Error != null && message.hasOwnProperty("Error"))
                if (!$util.isInteger(message.Error))
                    return "Error: integer expected";
            if (message.SocketId != null && message.hasOwnProperty("SocketId"))
                if (!$util.isInteger(message.SocketId))
                    return "SocketId: integer expected";
            return null;
        };

        /**
         * Creates a A_C_RegisterResponse message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof message.A_C_RegisterResponse
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {message.A_C_RegisterResponse} A_C_RegisterResponse
         */
        A_C_RegisterResponse.fromObject = function fromObject(object) {
            if (object instanceof $root.message.A_C_RegisterResponse)
                return object;
            var message = new $root.message.A_C_RegisterResponse();
            if (object.PacketHead != null) {
                if (typeof object.PacketHead !== "object")
                    throw TypeError(".message.A_C_RegisterResponse.PacketHead: object expected");
                message.PacketHead = $root.message.Ipacket.fromObject(object.PacketHead);
            }
            if (object.Error != null)
                message.Error = object.Error | 0;
            if (object.SocketId != null)
                message.SocketId = object.SocketId | 0;
            return message;
        };

        /**
         * Creates a plain object from a A_C_RegisterResponse message. Also converts values to other types if specified.
         * @function toObject
         * @memberof message.A_C_RegisterResponse
         * @static
         * @param {message.A_C_RegisterResponse} message A_C_RegisterResponse
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        A_C_RegisterResponse.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.defaults) {
                object.PacketHead = null;
                object.Error = 0;
                object.SocketId = 0;
            }
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                object.PacketHead = $root.message.Ipacket.toObject(message.PacketHead, options);
            if (message.Error != null && message.hasOwnProperty("Error"))
                object.Error = message.Error;
            if (message.SocketId != null && message.hasOwnProperty("SocketId"))
                object.SocketId = message.SocketId;
            return object;
        };

        /**
         * Converts this A_C_RegisterResponse to JSON.
         * @function toJSON
         * @memberof message.A_C_RegisterResponse
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        A_C_RegisterResponse.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return A_C_RegisterResponse;
    })();

    message.A_C_LoginRequest = (function() {

        /**
         * Properties of a A_C_LoginRequest.
         * @memberof message
         * @interface IA_C_LoginRequest
         * @property {message.IIpacket|null} [PacketHead] A_C_LoginRequest PacketHead
         * @property {number|null} [Error] A_C_LoginRequest Error
         * @property {number|null} [SocketId] A_C_LoginRequest SocketId
         * @property {string|null} [AccountName] A_C_LoginRequest AccountName
         */

        /**
         * Constructs a new A_C_LoginRequest.
         * @memberof message
         * @classdesc Represents a A_C_LoginRequest.
         * @implements IA_C_LoginRequest
         * @constructor
         * @param {message.IA_C_LoginRequest=} [properties] Properties to set
         */
        function A_C_LoginRequest(properties) {
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * A_C_LoginRequest PacketHead.
         * @member {message.IIpacket|null|undefined} PacketHead
         * @memberof message.A_C_LoginRequest
         * @instance
         */
        A_C_LoginRequest.prototype.PacketHead = null;

        /**
         * A_C_LoginRequest Error.
         * @member {number} Error
         * @memberof message.A_C_LoginRequest
         * @instance
         */
        A_C_LoginRequest.prototype.Error = 0;

        /**
         * A_C_LoginRequest SocketId.
         * @member {number} SocketId
         * @memberof message.A_C_LoginRequest
         * @instance
         */
        A_C_LoginRequest.prototype.SocketId = 0;

        /**
         * A_C_LoginRequest AccountName.
         * @member {string} AccountName
         * @memberof message.A_C_LoginRequest
         * @instance
         */
        A_C_LoginRequest.prototype.AccountName = "";

        /**
         * Creates a new A_C_LoginRequest instance using the specified properties.
         * @function create
         * @memberof message.A_C_LoginRequest
         * @static
         * @param {message.IA_C_LoginRequest=} [properties] Properties to set
         * @returns {message.A_C_LoginRequest} A_C_LoginRequest instance
         */
        A_C_LoginRequest.create = function create(properties) {
            return new A_C_LoginRequest(properties);
        };

        /**
         * Encodes the specified A_C_LoginRequest message. Does not implicitly {@link message.A_C_LoginRequest.verify|verify} messages.
         * @function encode
         * @memberof message.A_C_LoginRequest
         * @static
         * @param {message.IA_C_LoginRequest} message A_C_LoginRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        A_C_LoginRequest.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                $root.message.Ipacket.encode(message.PacketHead, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
            if (message.Error != null && message.hasOwnProperty("Error"))
                writer.uint32(/* id 2, wireType 0 =*/16).int32(message.Error);
            if (message.SocketId != null && message.hasOwnProperty("SocketId"))
                writer.uint32(/* id 3, wireType 0 =*/24).int32(message.SocketId);
            if (message.AccountName != null && message.hasOwnProperty("AccountName"))
                writer.uint32(/* id 4, wireType 2 =*/34).string(message.AccountName);
            return writer;
        };

        /**
         * Encodes the specified A_C_LoginRequest message, length delimited. Does not implicitly {@link message.A_C_LoginRequest.verify|verify} messages.
         * @function encodeDelimited
         * @memberof message.A_C_LoginRequest
         * @static
         * @param {message.IA_C_LoginRequest} message A_C_LoginRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        A_C_LoginRequest.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a A_C_LoginRequest message from the specified reader or buffer.
         * @function decode
         * @memberof message.A_C_LoginRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {message.A_C_LoginRequest} A_C_LoginRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        A_C_LoginRequest.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.A_C_LoginRequest();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.PacketHead = $root.message.Ipacket.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.Error = reader.int32();
                    break;
                case 3:
                    message.SocketId = reader.int32();
                    break;
                case 4:
                    message.AccountName = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a A_C_LoginRequest message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof message.A_C_LoginRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {message.A_C_LoginRequest} A_C_LoginRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        A_C_LoginRequest.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a A_C_LoginRequest message.
         * @function verify
         * @memberof message.A_C_LoginRequest
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        A_C_LoginRequest.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead")) {
                var error = $root.message.Ipacket.verify(message.PacketHead);
                if (error)
                    return "PacketHead." + error;
            }
            if (message.Error != null && message.hasOwnProperty("Error"))
                if (!$util.isInteger(message.Error))
                    return "Error: integer expected";
            if (message.SocketId != null && message.hasOwnProperty("SocketId"))
                if (!$util.isInteger(message.SocketId))
                    return "SocketId: integer expected";
            if (message.AccountName != null && message.hasOwnProperty("AccountName"))
                if (!$util.isString(message.AccountName))
                    return "AccountName: string expected";
            return null;
        };

        /**
         * Creates a A_C_LoginRequest message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof message.A_C_LoginRequest
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {message.A_C_LoginRequest} A_C_LoginRequest
         */
        A_C_LoginRequest.fromObject = function fromObject(object) {
            if (object instanceof $root.message.A_C_LoginRequest)
                return object;
            var message = new $root.message.A_C_LoginRequest();
            if (object.PacketHead != null) {
                if (typeof object.PacketHead !== "object")
                    throw TypeError(".message.A_C_LoginRequest.PacketHead: object expected");
                message.PacketHead = $root.message.Ipacket.fromObject(object.PacketHead);
            }
            if (object.Error != null)
                message.Error = object.Error | 0;
            if (object.SocketId != null)
                message.SocketId = object.SocketId | 0;
            if (object.AccountName != null)
                message.AccountName = String(object.AccountName);
            return message;
        };

        /**
         * Creates a plain object from a A_C_LoginRequest message. Also converts values to other types if specified.
         * @function toObject
         * @memberof message.A_C_LoginRequest
         * @static
         * @param {message.A_C_LoginRequest} message A_C_LoginRequest
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        A_C_LoginRequest.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.defaults) {
                object.PacketHead = null;
                object.Error = 0;
                object.SocketId = 0;
                object.AccountName = "";
            }
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                object.PacketHead = $root.message.Ipacket.toObject(message.PacketHead, options);
            if (message.Error != null && message.hasOwnProperty("Error"))
                object.Error = message.Error;
            if (message.SocketId != null && message.hasOwnProperty("SocketId"))
                object.SocketId = message.SocketId;
            if (message.AccountName != null && message.hasOwnProperty("AccountName"))
                object.AccountName = message.AccountName;
            return object;
        };

        /**
         * Converts this A_C_LoginRequest to JSON.
         * @function toJSON
         * @memberof message.A_C_LoginRequest
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        A_C_LoginRequest.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return A_C_LoginRequest;
    })();

    return message;
})();

module.exports = $root;
