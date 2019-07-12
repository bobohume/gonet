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

    message.C_A_LoginRequest = (function() {

        /**
         * Properties of a C_A_LoginRequest.
         * @memberof message
         * @interface IC_A_LoginRequest
         * @property {message.IIpacket|null} [PacketHead] C_A_LoginRequest PacketHead
         * @property {string|null} [AccountName] C_A_LoginRequest AccountName
         * @property {string|null} [BuildNo] C_A_LoginRequest BuildNo
         * @property {number|null} [SocketId] C_A_LoginRequest SocketId
         */

        /**
         * Constructs a new C_A_LoginRequest.
         * @memberof message
         * @classdesc Represents a C_A_LoginRequest.
         * @implements IC_A_LoginRequest
         * @constructor
         * @param {message.IC_A_LoginRequest=} [properties] Properties to set
         */
        function C_A_LoginRequest(properties) {
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * C_A_LoginRequest PacketHead.
         * @member {message.IIpacket|null|undefined} PacketHead
         * @memberof message.C_A_LoginRequest
         * @instance
         */
        C_A_LoginRequest.prototype.PacketHead = null;

        /**
         * C_A_LoginRequest AccountName.
         * @member {string} AccountName
         * @memberof message.C_A_LoginRequest
         * @instance
         */
        C_A_LoginRequest.prototype.AccountName = "";

        /**
         * C_A_LoginRequest BuildNo.
         * @member {string} BuildNo
         * @memberof message.C_A_LoginRequest
         * @instance
         */
        C_A_LoginRequest.prototype.BuildNo = "";

        /**
         * C_A_LoginRequest SocketId.
         * @member {number} SocketId
         * @memberof message.C_A_LoginRequest
         * @instance
         */
        C_A_LoginRequest.prototype.SocketId = 0;

        /**
         * Creates a new C_A_LoginRequest instance using the specified properties.
         * @function create
         * @memberof message.C_A_LoginRequest
         * @static
         * @param {message.IC_A_LoginRequest=} [properties] Properties to set
         * @returns {message.C_A_LoginRequest} C_A_LoginRequest instance
         */
        C_A_LoginRequest.create = function create(properties) {
            return new C_A_LoginRequest(properties);
        };

        /**
         * Encodes the specified C_A_LoginRequest message. Does not implicitly {@link message.C_A_LoginRequest.verify|verify} messages.
         * @function encode
         * @memberof message.C_A_LoginRequest
         * @static
         * @param {message.IC_A_LoginRequest} message C_A_LoginRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        C_A_LoginRequest.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                $root.message.Ipacket.encode(message.PacketHead, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
            if (message.AccountName != null && message.hasOwnProperty("AccountName"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.AccountName);
            if (message.BuildNo != null && message.hasOwnProperty("BuildNo"))
                writer.uint32(/* id 3, wireType 2 =*/26).string(message.BuildNo);
            if (message.SocketId != null && message.hasOwnProperty("SocketId"))
                writer.uint32(/* id 4, wireType 0 =*/32).int32(message.SocketId);
            return writer;
        };

        /**
         * Encodes the specified C_A_LoginRequest message, length delimited. Does not implicitly {@link message.C_A_LoginRequest.verify|verify} messages.
         * @function encodeDelimited
         * @memberof message.C_A_LoginRequest
         * @static
         * @param {message.IC_A_LoginRequest} message C_A_LoginRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        C_A_LoginRequest.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a C_A_LoginRequest message from the specified reader or buffer.
         * @function decode
         * @memberof message.C_A_LoginRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {message.C_A_LoginRequest} C_A_LoginRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        C_A_LoginRequest.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.C_A_LoginRequest();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.PacketHead = $root.message.Ipacket.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.AccountName = reader.string();
                    break;
                case 3:
                    message.BuildNo = reader.string();
                    break;
                case 4:
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
         * Decodes a C_A_LoginRequest message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof message.C_A_LoginRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {message.C_A_LoginRequest} C_A_LoginRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        C_A_LoginRequest.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a C_A_LoginRequest message.
         * @function verify
         * @memberof message.C_A_LoginRequest
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        C_A_LoginRequest.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead")) {
                var error = $root.message.Ipacket.verify(message.PacketHead);
                if (error)
                    return "PacketHead." + error;
            }
            if (message.AccountName != null && message.hasOwnProperty("AccountName"))
                if (!$util.isString(message.AccountName))
                    return "AccountName: string expected";
            if (message.BuildNo != null && message.hasOwnProperty("BuildNo"))
                if (!$util.isString(message.BuildNo))
                    return "BuildNo: string expected";
            if (message.SocketId != null && message.hasOwnProperty("SocketId"))
                if (!$util.isInteger(message.SocketId))
                    return "SocketId: integer expected";
            return null;
        };

        /**
         * Creates a C_A_LoginRequest message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof message.C_A_LoginRequest
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {message.C_A_LoginRequest} C_A_LoginRequest
         */
        C_A_LoginRequest.fromObject = function fromObject(object) {
            if (object instanceof $root.message.C_A_LoginRequest)
                return object;
            var message = new $root.message.C_A_LoginRequest();
            if (object.PacketHead != null) {
                if (typeof object.PacketHead !== "object")
                    throw TypeError(".message.C_A_LoginRequest.PacketHead: object expected");
                message.PacketHead = $root.message.Ipacket.fromObject(object.PacketHead);
            }
            if (object.AccountName != null)
                message.AccountName = String(object.AccountName);
            if (object.BuildNo != null)
                message.BuildNo = String(object.BuildNo);
            if (object.SocketId != null)
                message.SocketId = object.SocketId | 0;
            return message;
        };

        /**
         * Creates a plain object from a C_A_LoginRequest message. Also converts values to other types if specified.
         * @function toObject
         * @memberof message.C_A_LoginRequest
         * @static
         * @param {message.C_A_LoginRequest} message C_A_LoginRequest
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        C_A_LoginRequest.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.defaults) {
                object.PacketHead = null;
                object.AccountName = "";
                object.BuildNo = "";
                object.SocketId = 0;
            }
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                object.PacketHead = $root.message.Ipacket.toObject(message.PacketHead, options);
            if (message.AccountName != null && message.hasOwnProperty("AccountName"))
                object.AccountName = message.AccountName;
            if (message.BuildNo != null && message.hasOwnProperty("BuildNo"))
                object.BuildNo = message.BuildNo;
            if (message.SocketId != null && message.hasOwnProperty("SocketId"))
                object.SocketId = message.SocketId;
            return object;
        };

        /**
         * Converts this C_A_LoginRequest to JSON.
         * @function toJSON
         * @memberof message.C_A_LoginRequest
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        C_A_LoginRequest.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return C_A_LoginRequest;
    })();

    message.C_A_LoginRequest1 = (function() {

        /**
         * Properties of a C_A_LoginRequest1.
         * @memberof message
         * @interface IC_A_LoginRequest1
         * @property {message.IC_A_LoginRequest|null} [Login] C_A_LoginRequest1 Login
         */

        /**
         * Constructs a new C_A_LoginRequest1.
         * @memberof message
         * @classdesc Represents a C_A_LoginRequest1.
         * @implements IC_A_LoginRequest1
         * @constructor
         * @param {message.IC_A_LoginRequest1=} [properties] Properties to set
         */
        function C_A_LoginRequest1(properties) {
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * C_A_LoginRequest1 Login.
         * @member {message.IC_A_LoginRequest|null|undefined} Login
         * @memberof message.C_A_LoginRequest1
         * @instance
         */
        C_A_LoginRequest1.prototype.Login = null;

        /**
         * Creates a new C_A_LoginRequest1 instance using the specified properties.
         * @function create
         * @memberof message.C_A_LoginRequest1
         * @static
         * @param {message.IC_A_LoginRequest1=} [properties] Properties to set
         * @returns {message.C_A_LoginRequest1} C_A_LoginRequest1 instance
         */
        C_A_LoginRequest1.create = function create(properties) {
            return new C_A_LoginRequest1(properties);
        };

        /**
         * Encodes the specified C_A_LoginRequest1 message. Does not implicitly {@link message.C_A_LoginRequest1.verify|verify} messages.
         * @function encode
         * @memberof message.C_A_LoginRequest1
         * @static
         * @param {message.IC_A_LoginRequest1} message C_A_LoginRequest1 message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        C_A_LoginRequest1.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.Login != null && message.hasOwnProperty("Login"))
                $root.message.C_A_LoginRequest.encode(message.Login, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
            return writer;
        };

        /**
         * Encodes the specified C_A_LoginRequest1 message, length delimited. Does not implicitly {@link message.C_A_LoginRequest1.verify|verify} messages.
         * @function encodeDelimited
         * @memberof message.C_A_LoginRequest1
         * @static
         * @param {message.IC_A_LoginRequest1} message C_A_LoginRequest1 message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        C_A_LoginRequest1.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a C_A_LoginRequest1 message from the specified reader or buffer.
         * @function decode
         * @memberof message.C_A_LoginRequest1
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {message.C_A_LoginRequest1} C_A_LoginRequest1
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        C_A_LoginRequest1.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.C_A_LoginRequest1();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.Login = $root.message.C_A_LoginRequest.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a C_A_LoginRequest1 message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof message.C_A_LoginRequest1
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {message.C_A_LoginRequest1} C_A_LoginRequest1
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        C_A_LoginRequest1.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a C_A_LoginRequest1 message.
         * @function verify
         * @memberof message.C_A_LoginRequest1
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        C_A_LoginRequest1.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.Login != null && message.hasOwnProperty("Login")) {
                var error = $root.message.C_A_LoginRequest.verify(message.Login);
                if (error)
                    return "Login." + error;
            }
            return null;
        };

        /**
         * Creates a C_A_LoginRequest1 message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof message.C_A_LoginRequest1
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {message.C_A_LoginRequest1} C_A_LoginRequest1
         */
        C_A_LoginRequest1.fromObject = function fromObject(object) {
            if (object instanceof $root.message.C_A_LoginRequest1)
                return object;
            var message = new $root.message.C_A_LoginRequest1();
            if (object.Login != null) {
                if (typeof object.Login !== "object")
                    throw TypeError(".message.C_A_LoginRequest1.Login: object expected");
                message.Login = $root.message.C_A_LoginRequest.fromObject(object.Login);
            }
            return message;
        };

        /**
         * Creates a plain object from a C_A_LoginRequest1 message. Also converts values to other types if specified.
         * @function toObject
         * @memberof message.C_A_LoginRequest1
         * @static
         * @param {message.C_A_LoginRequest1} message C_A_LoginRequest1
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        C_A_LoginRequest1.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.defaults)
                object.Login = null;
            if (message.Login != null && message.hasOwnProperty("Login"))
                object.Login = $root.message.C_A_LoginRequest.toObject(message.Login, options);
            return object;
        };

        /**
         * Converts this C_A_LoginRequest1 to JSON.
         * @function toJSON
         * @memberof message.C_A_LoginRequest1
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        C_A_LoginRequest1.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return C_A_LoginRequest1;
    })();

    message.C_A_RegisterRequest = (function() {

        /**
         * Properties of a C_A_RegisterRequest.
         * @memberof message
         * @interface IC_A_RegisterRequest
         * @property {message.IIpacket|null} [PacketHead] C_A_RegisterRequest PacketHead
         * @property {string|null} [AccountName] C_A_RegisterRequest AccountName
         * @property {number|null} [SocketId] C_A_RegisterRequest SocketId
         */

        /**
         * Constructs a new C_A_RegisterRequest.
         * @memberof message
         * @classdesc Represents a C_A_RegisterRequest.
         * @implements IC_A_RegisterRequest
         * @constructor
         * @param {message.IC_A_RegisterRequest=} [properties] Properties to set
         */
        function C_A_RegisterRequest(properties) {
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * C_A_RegisterRequest PacketHead.
         * @member {message.IIpacket|null|undefined} PacketHead
         * @memberof message.C_A_RegisterRequest
         * @instance
         */
        C_A_RegisterRequest.prototype.PacketHead = null;

        /**
         * C_A_RegisterRequest AccountName.
         * @member {string} AccountName
         * @memberof message.C_A_RegisterRequest
         * @instance
         */
        C_A_RegisterRequest.prototype.AccountName = "";

        /**
         * C_A_RegisterRequest SocketId.
         * @member {number} SocketId
         * @memberof message.C_A_RegisterRequest
         * @instance
         */
        C_A_RegisterRequest.prototype.SocketId = 0;

        /**
         * Creates a new C_A_RegisterRequest instance using the specified properties.
         * @function create
         * @memberof message.C_A_RegisterRequest
         * @static
         * @param {message.IC_A_RegisterRequest=} [properties] Properties to set
         * @returns {message.C_A_RegisterRequest} C_A_RegisterRequest instance
         */
        C_A_RegisterRequest.create = function create(properties) {
            return new C_A_RegisterRequest(properties);
        };

        /**
         * Encodes the specified C_A_RegisterRequest message. Does not implicitly {@link message.C_A_RegisterRequest.verify|verify} messages.
         * @function encode
         * @memberof message.C_A_RegisterRequest
         * @static
         * @param {message.IC_A_RegisterRequest} message C_A_RegisterRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        C_A_RegisterRequest.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                $root.message.Ipacket.encode(message.PacketHead, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
            if (message.AccountName != null && message.hasOwnProperty("AccountName"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.AccountName);
            if (message.SocketId != null && message.hasOwnProperty("SocketId"))
                writer.uint32(/* id 3, wireType 0 =*/24).int32(message.SocketId);
            return writer;
        };

        /**
         * Encodes the specified C_A_RegisterRequest message, length delimited. Does not implicitly {@link message.C_A_RegisterRequest.verify|verify} messages.
         * @function encodeDelimited
         * @memberof message.C_A_RegisterRequest
         * @static
         * @param {message.IC_A_RegisterRequest} message C_A_RegisterRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        C_A_RegisterRequest.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a C_A_RegisterRequest message from the specified reader or buffer.
         * @function decode
         * @memberof message.C_A_RegisterRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {message.C_A_RegisterRequest} C_A_RegisterRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        C_A_RegisterRequest.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.C_A_RegisterRequest();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.PacketHead = $root.message.Ipacket.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.AccountName = reader.string();
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
         * Decodes a C_A_RegisterRequest message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof message.C_A_RegisterRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {message.C_A_RegisterRequest} C_A_RegisterRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        C_A_RegisterRequest.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a C_A_RegisterRequest message.
         * @function verify
         * @memberof message.C_A_RegisterRequest
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        C_A_RegisterRequest.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead")) {
                var error = $root.message.Ipacket.verify(message.PacketHead);
                if (error)
                    return "PacketHead." + error;
            }
            if (message.AccountName != null && message.hasOwnProperty("AccountName"))
                if (!$util.isString(message.AccountName))
                    return "AccountName: string expected";
            if (message.SocketId != null && message.hasOwnProperty("SocketId"))
                if (!$util.isInteger(message.SocketId))
                    return "SocketId: integer expected";
            return null;
        };

        /**
         * Creates a C_A_RegisterRequest message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof message.C_A_RegisterRequest
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {message.C_A_RegisterRequest} C_A_RegisterRequest
         */
        C_A_RegisterRequest.fromObject = function fromObject(object) {
            if (object instanceof $root.message.C_A_RegisterRequest)
                return object;
            var message = new $root.message.C_A_RegisterRequest();
            if (object.PacketHead != null) {
                if (typeof object.PacketHead !== "object")
                    throw TypeError(".message.C_A_RegisterRequest.PacketHead: object expected");
                message.PacketHead = $root.message.Ipacket.fromObject(object.PacketHead);
            }
            if (object.AccountName != null)
                message.AccountName = String(object.AccountName);
            if (object.SocketId != null)
                message.SocketId = object.SocketId | 0;
            return message;
        };

        /**
         * Creates a plain object from a C_A_RegisterRequest message. Also converts values to other types if specified.
         * @function toObject
         * @memberof message.C_A_RegisterRequest
         * @static
         * @param {message.C_A_RegisterRequest} message C_A_RegisterRequest
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        C_A_RegisterRequest.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.defaults) {
                object.PacketHead = null;
                object.AccountName = "";
                object.SocketId = 0;
            }
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                object.PacketHead = $root.message.Ipacket.toObject(message.PacketHead, options);
            if (message.AccountName != null && message.hasOwnProperty("AccountName"))
                object.AccountName = message.AccountName;
            if (message.SocketId != null && message.hasOwnProperty("SocketId"))
                object.SocketId = message.SocketId;
            return object;
        };

        /**
         * Converts this C_A_RegisterRequest to JSON.
         * @function toJSON
         * @memberof message.C_A_RegisterRequest
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        C_A_RegisterRequest.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return C_A_RegisterRequest;
    })();

    message.C_W_CreatePlayerRequest = (function() {

        /**
         * Properties of a C_W_CreatePlayerRequest.
         * @memberof message
         * @interface IC_W_CreatePlayerRequest
         * @property {message.IIpacket|null} [PacketHead] C_W_CreatePlayerRequest PacketHead
         * @property {string|null} [PlayerName] C_W_CreatePlayerRequest PlayerName
         * @property {number|null} [Sex] C_W_CreatePlayerRequest Sex
         */

        /**
         * Constructs a new C_W_CreatePlayerRequest.
         * @memberof message
         * @classdesc Represents a C_W_CreatePlayerRequest.
         * @implements IC_W_CreatePlayerRequest
         * @constructor
         * @param {message.IC_W_CreatePlayerRequest=} [properties] Properties to set
         */
        function C_W_CreatePlayerRequest(properties) {
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * C_W_CreatePlayerRequest PacketHead.
         * @member {message.IIpacket|null|undefined} PacketHead
         * @memberof message.C_W_CreatePlayerRequest
         * @instance
         */
        C_W_CreatePlayerRequest.prototype.PacketHead = null;

        /**
         * C_W_CreatePlayerRequest PlayerName.
         * @member {string} PlayerName
         * @memberof message.C_W_CreatePlayerRequest
         * @instance
         */
        C_W_CreatePlayerRequest.prototype.PlayerName = "";

        /**
         * C_W_CreatePlayerRequest Sex.
         * @member {number} Sex
         * @memberof message.C_W_CreatePlayerRequest
         * @instance
         */
        C_W_CreatePlayerRequest.prototype.Sex = 0;

        /**
         * Creates a new C_W_CreatePlayerRequest instance using the specified properties.
         * @function create
         * @memberof message.C_W_CreatePlayerRequest
         * @static
         * @param {message.IC_W_CreatePlayerRequest=} [properties] Properties to set
         * @returns {message.C_W_CreatePlayerRequest} C_W_CreatePlayerRequest instance
         */
        C_W_CreatePlayerRequest.create = function create(properties) {
            return new C_W_CreatePlayerRequest(properties);
        };

        /**
         * Encodes the specified C_W_CreatePlayerRequest message. Does not implicitly {@link message.C_W_CreatePlayerRequest.verify|verify} messages.
         * @function encode
         * @memberof message.C_W_CreatePlayerRequest
         * @static
         * @param {message.IC_W_CreatePlayerRequest} message C_W_CreatePlayerRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        C_W_CreatePlayerRequest.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                $root.message.Ipacket.encode(message.PacketHead, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
            if (message.PlayerName != null && message.hasOwnProperty("PlayerName"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.PlayerName);
            if (message.Sex != null && message.hasOwnProperty("Sex"))
                writer.uint32(/* id 3, wireType 0 =*/24).int32(message.Sex);
            return writer;
        };

        /**
         * Encodes the specified C_W_CreatePlayerRequest message, length delimited. Does not implicitly {@link message.C_W_CreatePlayerRequest.verify|verify} messages.
         * @function encodeDelimited
         * @memberof message.C_W_CreatePlayerRequest
         * @static
         * @param {message.IC_W_CreatePlayerRequest} message C_W_CreatePlayerRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        C_W_CreatePlayerRequest.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a C_W_CreatePlayerRequest message from the specified reader or buffer.
         * @function decode
         * @memberof message.C_W_CreatePlayerRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {message.C_W_CreatePlayerRequest} C_W_CreatePlayerRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        C_W_CreatePlayerRequest.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.C_W_CreatePlayerRequest();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.PacketHead = $root.message.Ipacket.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.PlayerName = reader.string();
                    break;
                case 3:
                    message.Sex = reader.int32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a C_W_CreatePlayerRequest message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof message.C_W_CreatePlayerRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {message.C_W_CreatePlayerRequest} C_W_CreatePlayerRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        C_W_CreatePlayerRequest.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a C_W_CreatePlayerRequest message.
         * @function verify
         * @memberof message.C_W_CreatePlayerRequest
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        C_W_CreatePlayerRequest.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead")) {
                var error = $root.message.Ipacket.verify(message.PacketHead);
                if (error)
                    return "PacketHead." + error;
            }
            if (message.PlayerName != null && message.hasOwnProperty("PlayerName"))
                if (!$util.isString(message.PlayerName))
                    return "PlayerName: string expected";
            if (message.Sex != null && message.hasOwnProperty("Sex"))
                if (!$util.isInteger(message.Sex))
                    return "Sex: integer expected";
            return null;
        };

        /**
         * Creates a C_W_CreatePlayerRequest message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof message.C_W_CreatePlayerRequest
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {message.C_W_CreatePlayerRequest} C_W_CreatePlayerRequest
         */
        C_W_CreatePlayerRequest.fromObject = function fromObject(object) {
            if (object instanceof $root.message.C_W_CreatePlayerRequest)
                return object;
            var message = new $root.message.C_W_CreatePlayerRequest();
            if (object.PacketHead != null) {
                if (typeof object.PacketHead !== "object")
                    throw TypeError(".message.C_W_CreatePlayerRequest.PacketHead: object expected");
                message.PacketHead = $root.message.Ipacket.fromObject(object.PacketHead);
            }
            if (object.PlayerName != null)
                message.PlayerName = String(object.PlayerName);
            if (object.Sex != null)
                message.Sex = object.Sex | 0;
            return message;
        };

        /**
         * Creates a plain object from a C_W_CreatePlayerRequest message. Also converts values to other types if specified.
         * @function toObject
         * @memberof message.C_W_CreatePlayerRequest
         * @static
         * @param {message.C_W_CreatePlayerRequest} message C_W_CreatePlayerRequest
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        C_W_CreatePlayerRequest.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.defaults) {
                object.PacketHead = null;
                object.PlayerName = "";
                object.Sex = 0;
            }
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                object.PacketHead = $root.message.Ipacket.toObject(message.PacketHead, options);
            if (message.PlayerName != null && message.hasOwnProperty("PlayerName"))
                object.PlayerName = message.PlayerName;
            if (message.Sex != null && message.hasOwnProperty("Sex"))
                object.Sex = message.Sex;
            return object;
        };

        /**
         * Converts this C_W_CreatePlayerRequest to JSON.
         * @function toJSON
         * @memberof message.C_W_CreatePlayerRequest
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        C_W_CreatePlayerRequest.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return C_W_CreatePlayerRequest;
    })();

    message.C_W_Game_LoginRequset = (function() {

        /**
         * Properties of a C_W_Game_LoginRequset.
         * @memberof message
         * @interface IC_W_Game_LoginRequset
         * @property {message.IIpacket|null} [PacketHead] C_W_Game_LoginRequset PacketHead
         * @property {number|Long|null} [PlayerId] C_W_Game_LoginRequset PlayerId
         */

        /**
         * Constructs a new C_W_Game_LoginRequset.
         * @memberof message
         * @classdesc Represents a C_W_Game_LoginRequset.
         * @implements IC_W_Game_LoginRequset
         * @constructor
         * @param {message.IC_W_Game_LoginRequset=} [properties] Properties to set
         */
        function C_W_Game_LoginRequset(properties) {
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * C_W_Game_LoginRequset PacketHead.
         * @member {message.IIpacket|null|undefined} PacketHead
         * @memberof message.C_W_Game_LoginRequset
         * @instance
         */
        C_W_Game_LoginRequset.prototype.PacketHead = null;

        /**
         * C_W_Game_LoginRequset PlayerId.
         * @member {number|Long} PlayerId
         * @memberof message.C_W_Game_LoginRequset
         * @instance
         */
        C_W_Game_LoginRequset.prototype.PlayerId = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * Creates a new C_W_Game_LoginRequset instance using the specified properties.
         * @function create
         * @memberof message.C_W_Game_LoginRequset
         * @static
         * @param {message.IC_W_Game_LoginRequset=} [properties] Properties to set
         * @returns {message.C_W_Game_LoginRequset} C_W_Game_LoginRequset instance
         */
        C_W_Game_LoginRequset.create = function create(properties) {
            return new C_W_Game_LoginRequset(properties);
        };

        /**
         * Encodes the specified C_W_Game_LoginRequset message. Does not implicitly {@link message.C_W_Game_LoginRequset.verify|verify} messages.
         * @function encode
         * @memberof message.C_W_Game_LoginRequset
         * @static
         * @param {message.IC_W_Game_LoginRequset} message C_W_Game_LoginRequset message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        C_W_Game_LoginRequset.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                $root.message.Ipacket.encode(message.PacketHead, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
            if (message.PlayerId != null && message.hasOwnProperty("PlayerId"))
                writer.uint32(/* id 2, wireType 0 =*/16).int64(message.PlayerId);
            return writer;
        };

        /**
         * Encodes the specified C_W_Game_LoginRequset message, length delimited. Does not implicitly {@link message.C_W_Game_LoginRequset.verify|verify} messages.
         * @function encodeDelimited
         * @memberof message.C_W_Game_LoginRequset
         * @static
         * @param {message.IC_W_Game_LoginRequset} message C_W_Game_LoginRequset message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        C_W_Game_LoginRequset.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a C_W_Game_LoginRequset message from the specified reader or buffer.
         * @function decode
         * @memberof message.C_W_Game_LoginRequset
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {message.C_W_Game_LoginRequset} C_W_Game_LoginRequset
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        C_W_Game_LoginRequset.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.C_W_Game_LoginRequset();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.PacketHead = $root.message.Ipacket.decode(reader, reader.uint32());
                    break;
                case 2:
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
         * Decodes a C_W_Game_LoginRequset message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof message.C_W_Game_LoginRequset
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {message.C_W_Game_LoginRequset} C_W_Game_LoginRequset
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        C_W_Game_LoginRequset.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a C_W_Game_LoginRequset message.
         * @function verify
         * @memberof message.C_W_Game_LoginRequset
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        C_W_Game_LoginRequset.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead")) {
                var error = $root.message.Ipacket.verify(message.PacketHead);
                if (error)
                    return "PacketHead." + error;
            }
            if (message.PlayerId != null && message.hasOwnProperty("PlayerId"))
                if (!$util.isInteger(message.PlayerId) && !(message.PlayerId && $util.isInteger(message.PlayerId.low) && $util.isInteger(message.PlayerId.high)))
                    return "PlayerId: integer|Long expected";
            return null;
        };

        /**
         * Creates a C_W_Game_LoginRequset message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof message.C_W_Game_LoginRequset
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {message.C_W_Game_LoginRequset} C_W_Game_LoginRequset
         */
        C_W_Game_LoginRequset.fromObject = function fromObject(object) {
            if (object instanceof $root.message.C_W_Game_LoginRequset)
                return object;
            var message = new $root.message.C_W_Game_LoginRequset();
            if (object.PacketHead != null) {
                if (typeof object.PacketHead !== "object")
                    throw TypeError(".message.C_W_Game_LoginRequset.PacketHead: object expected");
                message.PacketHead = $root.message.Ipacket.fromObject(object.PacketHead);
            }
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
         * Creates a plain object from a C_W_Game_LoginRequset message. Also converts values to other types if specified.
         * @function toObject
         * @memberof message.C_W_Game_LoginRequset
         * @static
         * @param {message.C_W_Game_LoginRequset} message C_W_Game_LoginRequset
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        C_W_Game_LoginRequset.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.defaults) {
                object.PacketHead = null;
                if ($util.Long) {
                    var long = new $util.Long(0, 0, false);
                    object.PlayerId = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.PlayerId = options.longs === String ? "0" : 0;
            }
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                object.PacketHead = $root.message.Ipacket.toObject(message.PacketHead, options);
            if (message.PlayerId != null && message.hasOwnProperty("PlayerId"))
                if (typeof message.PlayerId === "number")
                    object.PlayerId = options.longs === String ? String(message.PlayerId) : message.PlayerId;
                else
                    object.PlayerId = options.longs === String ? $util.Long.prototype.toString.call(message.PlayerId) : options.longs === Number ? new $util.LongBits(message.PlayerId.low >>> 0, message.PlayerId.high >>> 0).toNumber() : message.PlayerId;
            return object;
        };

        /**
         * Converts this C_W_Game_LoginRequset to JSON.
         * @function toJSON
         * @memberof message.C_W_Game_LoginRequset
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        C_W_Game_LoginRequset.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return C_W_Game_LoginRequset;
    })();

    message.C_G_LogoutResponse = (function() {

        /**
         * Properties of a C_G_LogoutResponse.
         * @memberof message
         * @interface IC_G_LogoutResponse
         * @property {message.IIpacket|null} [PacketHead] C_G_LogoutResponse PacketHead
         */

        /**
         * Constructs a new C_G_LogoutResponse.
         * @memberof message
         * @classdesc Represents a C_G_LogoutResponse.
         * @implements IC_G_LogoutResponse
         * @constructor
         * @param {message.IC_G_LogoutResponse=} [properties] Properties to set
         */
        function C_G_LogoutResponse(properties) {
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * C_G_LogoutResponse PacketHead.
         * @member {message.IIpacket|null|undefined} PacketHead
         * @memberof message.C_G_LogoutResponse
         * @instance
         */
        C_G_LogoutResponse.prototype.PacketHead = null;

        /**
         * Creates a new C_G_LogoutResponse instance using the specified properties.
         * @function create
         * @memberof message.C_G_LogoutResponse
         * @static
         * @param {message.IC_G_LogoutResponse=} [properties] Properties to set
         * @returns {message.C_G_LogoutResponse} C_G_LogoutResponse instance
         */
        C_G_LogoutResponse.create = function create(properties) {
            return new C_G_LogoutResponse(properties);
        };

        /**
         * Encodes the specified C_G_LogoutResponse message. Does not implicitly {@link message.C_G_LogoutResponse.verify|verify} messages.
         * @function encode
         * @memberof message.C_G_LogoutResponse
         * @static
         * @param {message.IC_G_LogoutResponse} message C_G_LogoutResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        C_G_LogoutResponse.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                $root.message.Ipacket.encode(message.PacketHead, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
            return writer;
        };

        /**
         * Encodes the specified C_G_LogoutResponse message, length delimited. Does not implicitly {@link message.C_G_LogoutResponse.verify|verify} messages.
         * @function encodeDelimited
         * @memberof message.C_G_LogoutResponse
         * @static
         * @param {message.IC_G_LogoutResponse} message C_G_LogoutResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        C_G_LogoutResponse.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a C_G_LogoutResponse message from the specified reader or buffer.
         * @function decode
         * @memberof message.C_G_LogoutResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {message.C_G_LogoutResponse} C_G_LogoutResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        C_G_LogoutResponse.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.C_G_LogoutResponse();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.PacketHead = $root.message.Ipacket.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a C_G_LogoutResponse message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof message.C_G_LogoutResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {message.C_G_LogoutResponse} C_G_LogoutResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        C_G_LogoutResponse.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a C_G_LogoutResponse message.
         * @function verify
         * @memberof message.C_G_LogoutResponse
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        C_G_LogoutResponse.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead")) {
                var error = $root.message.Ipacket.verify(message.PacketHead);
                if (error)
                    return "PacketHead." + error;
            }
            return null;
        };

        /**
         * Creates a C_G_LogoutResponse message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof message.C_G_LogoutResponse
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {message.C_G_LogoutResponse} C_G_LogoutResponse
         */
        C_G_LogoutResponse.fromObject = function fromObject(object) {
            if (object instanceof $root.message.C_G_LogoutResponse)
                return object;
            var message = new $root.message.C_G_LogoutResponse();
            if (object.PacketHead != null) {
                if (typeof object.PacketHead !== "object")
                    throw TypeError(".message.C_G_LogoutResponse.PacketHead: object expected");
                message.PacketHead = $root.message.Ipacket.fromObject(object.PacketHead);
            }
            return message;
        };

        /**
         * Creates a plain object from a C_G_LogoutResponse message. Also converts values to other types if specified.
         * @function toObject
         * @memberof message.C_G_LogoutResponse
         * @static
         * @param {message.C_G_LogoutResponse} message C_G_LogoutResponse
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        C_G_LogoutResponse.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.defaults)
                object.PacketHead = null;
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                object.PacketHead = $root.message.Ipacket.toObject(message.PacketHead, options);
            return object;
        };

        /**
         * Converts this C_G_LogoutResponse to JSON.
         * @function toJSON
         * @memberof message.C_G_LogoutResponse
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        C_G_LogoutResponse.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return C_G_LogoutResponse;
    })();

    message.C_W_ChatMessage = (function() {

        /**
         * Properties of a C_W_ChatMessage.
         * @memberof message
         * @interface IC_W_ChatMessage
         * @property {message.IIpacket|null} [PacketHead] C_W_ChatMessage PacketHead
         * @property {number|Long|null} [Sender] C_W_ChatMessage Sender
         * @property {number|Long|null} [Recver] C_W_ChatMessage Recver
         * @property {number|null} [MessageType] C_W_ChatMessage MessageType
         * @property {string|null} [Message] C_W_ChatMessage Message
         */

        /**
         * Constructs a new C_W_ChatMessage.
         * @memberof message
         * @classdesc Represents a C_W_ChatMessage.
         * @implements IC_W_ChatMessage
         * @constructor
         * @param {message.IC_W_ChatMessage=} [properties] Properties to set
         */
        function C_W_ChatMessage(properties) {
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * C_W_ChatMessage PacketHead.
         * @member {message.IIpacket|null|undefined} PacketHead
         * @memberof message.C_W_ChatMessage
         * @instance
         */
        C_W_ChatMessage.prototype.PacketHead = null;

        /**
         * C_W_ChatMessage Sender.
         * @member {number|Long} Sender
         * @memberof message.C_W_ChatMessage
         * @instance
         */
        C_W_ChatMessage.prototype.Sender = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * C_W_ChatMessage Recver.
         * @member {number|Long} Recver
         * @memberof message.C_W_ChatMessage
         * @instance
         */
        C_W_ChatMessage.prototype.Recver = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * C_W_ChatMessage MessageType.
         * @member {number} MessageType
         * @memberof message.C_W_ChatMessage
         * @instance
         */
        C_W_ChatMessage.prototype.MessageType = 0;

        /**
         * C_W_ChatMessage Message.
         * @member {string} Message
         * @memberof message.C_W_ChatMessage
         * @instance
         */
        C_W_ChatMessage.prototype.Message = "";

        /**
         * Creates a new C_W_ChatMessage instance using the specified properties.
         * @function create
         * @memberof message.C_W_ChatMessage
         * @static
         * @param {message.IC_W_ChatMessage=} [properties] Properties to set
         * @returns {message.C_W_ChatMessage} C_W_ChatMessage instance
         */
        C_W_ChatMessage.create = function create(properties) {
            return new C_W_ChatMessage(properties);
        };

        /**
         * Encodes the specified C_W_ChatMessage message. Does not implicitly {@link message.C_W_ChatMessage.verify|verify} messages.
         * @function encode
         * @memberof message.C_W_ChatMessage
         * @static
         * @param {message.IC_W_ChatMessage} message C_W_ChatMessage message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        C_W_ChatMessage.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                $root.message.Ipacket.encode(message.PacketHead, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
            if (message.Sender != null && message.hasOwnProperty("Sender"))
                writer.uint32(/* id 2, wireType 0 =*/16).int64(message.Sender);
            if (message.Recver != null && message.hasOwnProperty("Recver"))
                writer.uint32(/* id 3, wireType 0 =*/24).int64(message.Recver);
            if (message.MessageType != null && message.hasOwnProperty("MessageType"))
                writer.uint32(/* id 4, wireType 0 =*/32).int32(message.MessageType);
            if (message.Message != null && message.hasOwnProperty("Message"))
                writer.uint32(/* id 5, wireType 2 =*/42).string(message.Message);
            return writer;
        };

        /**
         * Encodes the specified C_W_ChatMessage message, length delimited. Does not implicitly {@link message.C_W_ChatMessage.verify|verify} messages.
         * @function encodeDelimited
         * @memberof message.C_W_ChatMessage
         * @static
         * @param {message.IC_W_ChatMessage} message C_W_ChatMessage message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        C_W_ChatMessage.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a C_W_ChatMessage message from the specified reader or buffer.
         * @function decode
         * @memberof message.C_W_ChatMessage
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {message.C_W_ChatMessage} C_W_ChatMessage
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        C_W_ChatMessage.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.C_W_ChatMessage();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.PacketHead = $root.message.Ipacket.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.Sender = reader.int64();
                    break;
                case 3:
                    message.Recver = reader.int64();
                    break;
                case 4:
                    message.MessageType = reader.int32();
                    break;
                case 5:
                    message.Message = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a C_W_ChatMessage message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof message.C_W_ChatMessage
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {message.C_W_ChatMessage} C_W_ChatMessage
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        C_W_ChatMessage.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a C_W_ChatMessage message.
         * @function verify
         * @memberof message.C_W_ChatMessage
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        C_W_ChatMessage.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead")) {
                var error = $root.message.Ipacket.verify(message.PacketHead);
                if (error)
                    return "PacketHead." + error;
            }
            if (message.Sender != null && message.hasOwnProperty("Sender"))
                if (!$util.isInteger(message.Sender) && !(message.Sender && $util.isInteger(message.Sender.low) && $util.isInteger(message.Sender.high)))
                    return "Sender: integer|Long expected";
            if (message.Recver != null && message.hasOwnProperty("Recver"))
                if (!$util.isInteger(message.Recver) && !(message.Recver && $util.isInteger(message.Recver.low) && $util.isInteger(message.Recver.high)))
                    return "Recver: integer|Long expected";
            if (message.MessageType != null && message.hasOwnProperty("MessageType"))
                if (!$util.isInteger(message.MessageType))
                    return "MessageType: integer expected";
            if (message.Message != null && message.hasOwnProperty("Message"))
                if (!$util.isString(message.Message))
                    return "Message: string expected";
            return null;
        };

        /**
         * Creates a C_W_ChatMessage message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof message.C_W_ChatMessage
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {message.C_W_ChatMessage} C_W_ChatMessage
         */
        C_W_ChatMessage.fromObject = function fromObject(object) {
            if (object instanceof $root.message.C_W_ChatMessage)
                return object;
            var message = new $root.message.C_W_ChatMessage();
            if (object.PacketHead != null) {
                if (typeof object.PacketHead !== "object")
                    throw TypeError(".message.C_W_ChatMessage.PacketHead: object expected");
                message.PacketHead = $root.message.Ipacket.fromObject(object.PacketHead);
            }
            if (object.Sender != null)
                if ($util.Long)
                    (message.Sender = $util.Long.fromValue(object.Sender)).unsigned = false;
                else if (typeof object.Sender === "string")
                    message.Sender = parseInt(object.Sender, 10);
                else if (typeof object.Sender === "number")
                    message.Sender = object.Sender;
                else if (typeof object.Sender === "object")
                    message.Sender = new $util.LongBits(object.Sender.low >>> 0, object.Sender.high >>> 0).toNumber();
            if (object.Recver != null)
                if ($util.Long)
                    (message.Recver = $util.Long.fromValue(object.Recver)).unsigned = false;
                else if (typeof object.Recver === "string")
                    message.Recver = parseInt(object.Recver, 10);
                else if (typeof object.Recver === "number")
                    message.Recver = object.Recver;
                else if (typeof object.Recver === "object")
                    message.Recver = new $util.LongBits(object.Recver.low >>> 0, object.Recver.high >>> 0).toNumber();
            if (object.MessageType != null)
                message.MessageType = object.MessageType | 0;
            if (object.Message != null)
                message.Message = String(object.Message);
            return message;
        };

        /**
         * Creates a plain object from a C_W_ChatMessage message. Also converts values to other types if specified.
         * @function toObject
         * @memberof message.C_W_ChatMessage
         * @static
         * @param {message.C_W_ChatMessage} message C_W_ChatMessage
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        C_W_ChatMessage.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.defaults) {
                object.PacketHead = null;
                if ($util.Long) {
                    var long = new $util.Long(0, 0, false);
                    object.Sender = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.Sender = options.longs === String ? "0" : 0;
                if ($util.Long) {
                    var long = new $util.Long(0, 0, false);
                    object.Recver = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.Recver = options.longs === String ? "0" : 0;
                object.MessageType = 0;
                object.Message = "";
            }
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                object.PacketHead = $root.message.Ipacket.toObject(message.PacketHead, options);
            if (message.Sender != null && message.hasOwnProperty("Sender"))
                if (typeof message.Sender === "number")
                    object.Sender = options.longs === String ? String(message.Sender) : message.Sender;
                else
                    object.Sender = options.longs === String ? $util.Long.prototype.toString.call(message.Sender) : options.longs === Number ? new $util.LongBits(message.Sender.low >>> 0, message.Sender.high >>> 0).toNumber() : message.Sender;
            if (message.Recver != null && message.hasOwnProperty("Recver"))
                if (typeof message.Recver === "number")
                    object.Recver = options.longs === String ? String(message.Recver) : message.Recver;
                else
                    object.Recver = options.longs === String ? $util.Long.prototype.toString.call(message.Recver) : options.longs === Number ? new $util.LongBits(message.Recver.low >>> 0, message.Recver.high >>> 0).toNumber() : message.Recver;
            if (message.MessageType != null && message.hasOwnProperty("MessageType"))
                object.MessageType = message.MessageType;
            if (message.Message != null && message.hasOwnProperty("Message"))
                object.Message = message.Message;
            return object;
        };

        /**
         * Converts this C_W_ChatMessage to JSON.
         * @function toJSON
         * @memberof message.C_W_ChatMessage
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        C_W_ChatMessage.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return C_W_ChatMessage;
    })();

    message.W_C_ChatMessage = (function() {

        /**
         * Properties of a W_C_ChatMessage.
         * @memberof message
         * @interface IW_C_ChatMessage
         * @property {message.IIpacket|null} [PacketHead] W_C_ChatMessage PacketHead
         * @property {number|Long|null} [Sender] W_C_ChatMessage Sender
         * @property {string|null} [SenderName] W_C_ChatMessage SenderName
         * @property {number|Long|null} [Recver] W_C_ChatMessage Recver
         * @property {string|null} [RecverName] W_C_ChatMessage RecverName
         * @property {number|null} [MessageType] W_C_ChatMessage MessageType
         * @property {string|null} [Message] W_C_ChatMessage Message
         */

        /**
         * Constructs a new W_C_ChatMessage.
         * @memberof message
         * @classdesc Represents a W_C_ChatMessage.
         * @implements IW_C_ChatMessage
         * @constructor
         * @param {message.IW_C_ChatMessage=} [properties] Properties to set
         */
        function W_C_ChatMessage(properties) {
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * W_C_ChatMessage PacketHead.
         * @member {message.IIpacket|null|undefined} PacketHead
         * @memberof message.W_C_ChatMessage
         * @instance
         */
        W_C_ChatMessage.prototype.PacketHead = null;

        /**
         * W_C_ChatMessage Sender.
         * @member {number|Long} Sender
         * @memberof message.W_C_ChatMessage
         * @instance
         */
        W_C_ChatMessage.prototype.Sender = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * W_C_ChatMessage SenderName.
         * @member {string} SenderName
         * @memberof message.W_C_ChatMessage
         * @instance
         */
        W_C_ChatMessage.prototype.SenderName = "";

        /**
         * W_C_ChatMessage Recver.
         * @member {number|Long} Recver
         * @memberof message.W_C_ChatMessage
         * @instance
         */
        W_C_ChatMessage.prototype.Recver = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * W_C_ChatMessage RecverName.
         * @member {string} RecverName
         * @memberof message.W_C_ChatMessage
         * @instance
         */
        W_C_ChatMessage.prototype.RecverName = "";

        /**
         * W_C_ChatMessage MessageType.
         * @member {number} MessageType
         * @memberof message.W_C_ChatMessage
         * @instance
         */
        W_C_ChatMessage.prototype.MessageType = 0;

        /**
         * W_C_ChatMessage Message.
         * @member {string} Message
         * @memberof message.W_C_ChatMessage
         * @instance
         */
        W_C_ChatMessage.prototype.Message = "";

        /**
         * Creates a new W_C_ChatMessage instance using the specified properties.
         * @function create
         * @memberof message.W_C_ChatMessage
         * @static
         * @param {message.IW_C_ChatMessage=} [properties] Properties to set
         * @returns {message.W_C_ChatMessage} W_C_ChatMessage instance
         */
        W_C_ChatMessage.create = function create(properties) {
            return new W_C_ChatMessage(properties);
        };

        /**
         * Encodes the specified W_C_ChatMessage message. Does not implicitly {@link message.W_C_ChatMessage.verify|verify} messages.
         * @function encode
         * @memberof message.W_C_ChatMessage
         * @static
         * @param {message.IW_C_ChatMessage} message W_C_ChatMessage message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        W_C_ChatMessage.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                $root.message.Ipacket.encode(message.PacketHead, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
            if (message.Sender != null && message.hasOwnProperty("Sender"))
                writer.uint32(/* id 2, wireType 0 =*/16).int64(message.Sender);
            if (message.SenderName != null && message.hasOwnProperty("SenderName"))
                writer.uint32(/* id 3, wireType 2 =*/26).string(message.SenderName);
            if (message.Recver != null && message.hasOwnProperty("Recver"))
                writer.uint32(/* id 4, wireType 0 =*/32).int64(message.Recver);
            if (message.RecverName != null && message.hasOwnProperty("RecverName"))
                writer.uint32(/* id 5, wireType 2 =*/42).string(message.RecverName);
            if (message.MessageType != null && message.hasOwnProperty("MessageType"))
                writer.uint32(/* id 6, wireType 0 =*/48).int32(message.MessageType);
            if (message.Message != null && message.hasOwnProperty("Message"))
                writer.uint32(/* id 7, wireType 2 =*/58).string(message.Message);
            return writer;
        };

        /**
         * Encodes the specified W_C_ChatMessage message, length delimited. Does not implicitly {@link message.W_C_ChatMessage.verify|verify} messages.
         * @function encodeDelimited
         * @memberof message.W_C_ChatMessage
         * @static
         * @param {message.IW_C_ChatMessage} message W_C_ChatMessage message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        W_C_ChatMessage.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a W_C_ChatMessage message from the specified reader or buffer.
         * @function decode
         * @memberof message.W_C_ChatMessage
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {message.W_C_ChatMessage} W_C_ChatMessage
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        W_C_ChatMessage.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.W_C_ChatMessage();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.PacketHead = $root.message.Ipacket.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.Sender = reader.int64();
                    break;
                case 3:
                    message.SenderName = reader.string();
                    break;
                case 4:
                    message.Recver = reader.int64();
                    break;
                case 5:
                    message.RecverName = reader.string();
                    break;
                case 6:
                    message.MessageType = reader.int32();
                    break;
                case 7:
                    message.Message = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a W_C_ChatMessage message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof message.W_C_ChatMessage
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {message.W_C_ChatMessage} W_C_ChatMessage
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        W_C_ChatMessage.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a W_C_ChatMessage message.
         * @function verify
         * @memberof message.W_C_ChatMessage
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        W_C_ChatMessage.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead")) {
                var error = $root.message.Ipacket.verify(message.PacketHead);
                if (error)
                    return "PacketHead." + error;
            }
            if (message.Sender != null && message.hasOwnProperty("Sender"))
                if (!$util.isInteger(message.Sender) && !(message.Sender && $util.isInteger(message.Sender.low) && $util.isInteger(message.Sender.high)))
                    return "Sender: integer|Long expected";
            if (message.SenderName != null && message.hasOwnProperty("SenderName"))
                if (!$util.isString(message.SenderName))
                    return "SenderName: string expected";
            if (message.Recver != null && message.hasOwnProperty("Recver"))
                if (!$util.isInteger(message.Recver) && !(message.Recver && $util.isInteger(message.Recver.low) && $util.isInteger(message.Recver.high)))
                    return "Recver: integer|Long expected";
            if (message.RecverName != null && message.hasOwnProperty("RecverName"))
                if (!$util.isString(message.RecverName))
                    return "RecverName: string expected";
            if (message.MessageType != null && message.hasOwnProperty("MessageType"))
                if (!$util.isInteger(message.MessageType))
                    return "MessageType: integer expected";
            if (message.Message != null && message.hasOwnProperty("Message"))
                if (!$util.isString(message.Message))
                    return "Message: string expected";
            return null;
        };

        /**
         * Creates a W_C_ChatMessage message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof message.W_C_ChatMessage
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {message.W_C_ChatMessage} W_C_ChatMessage
         */
        W_C_ChatMessage.fromObject = function fromObject(object) {
            if (object instanceof $root.message.W_C_ChatMessage)
                return object;
            var message = new $root.message.W_C_ChatMessage();
            if (object.PacketHead != null) {
                if (typeof object.PacketHead !== "object")
                    throw TypeError(".message.W_C_ChatMessage.PacketHead: object expected");
                message.PacketHead = $root.message.Ipacket.fromObject(object.PacketHead);
            }
            if (object.Sender != null)
                if ($util.Long)
                    (message.Sender = $util.Long.fromValue(object.Sender)).unsigned = false;
                else if (typeof object.Sender === "string")
                    message.Sender = parseInt(object.Sender, 10);
                else if (typeof object.Sender === "number")
                    message.Sender = object.Sender;
                else if (typeof object.Sender === "object")
                    message.Sender = new $util.LongBits(object.Sender.low >>> 0, object.Sender.high >>> 0).toNumber();
            if (object.SenderName != null)
                message.SenderName = String(object.SenderName);
            if (object.Recver != null)
                if ($util.Long)
                    (message.Recver = $util.Long.fromValue(object.Recver)).unsigned = false;
                else if (typeof object.Recver === "string")
                    message.Recver = parseInt(object.Recver, 10);
                else if (typeof object.Recver === "number")
                    message.Recver = object.Recver;
                else if (typeof object.Recver === "object")
                    message.Recver = new $util.LongBits(object.Recver.low >>> 0, object.Recver.high >>> 0).toNumber();
            if (object.RecverName != null)
                message.RecverName = String(object.RecverName);
            if (object.MessageType != null)
                message.MessageType = object.MessageType | 0;
            if (object.Message != null)
                message.Message = String(object.Message);
            return message;
        };

        /**
         * Creates a plain object from a W_C_ChatMessage message. Also converts values to other types if specified.
         * @function toObject
         * @memberof message.W_C_ChatMessage
         * @static
         * @param {message.W_C_ChatMessage} message W_C_ChatMessage
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        W_C_ChatMessage.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.defaults) {
                object.PacketHead = null;
                if ($util.Long) {
                    var long = new $util.Long(0, 0, false);
                    object.Sender = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.Sender = options.longs === String ? "0" : 0;
                object.SenderName = "";
                if ($util.Long) {
                    var long = new $util.Long(0, 0, false);
                    object.Recver = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.Recver = options.longs === String ? "0" : 0;
                object.RecverName = "";
                object.MessageType = 0;
                object.Message = "";
            }
            if (message.PacketHead != null && message.hasOwnProperty("PacketHead"))
                object.PacketHead = $root.message.Ipacket.toObject(message.PacketHead, options);
            if (message.Sender != null && message.hasOwnProperty("Sender"))
                if (typeof message.Sender === "number")
                    object.Sender = options.longs === String ? String(message.Sender) : message.Sender;
                else
                    object.Sender = options.longs === String ? $util.Long.prototype.toString.call(message.Sender) : options.longs === Number ? new $util.LongBits(message.Sender.low >>> 0, message.Sender.high >>> 0).toNumber() : message.Sender;
            if (message.SenderName != null && message.hasOwnProperty("SenderName"))
                object.SenderName = message.SenderName;
            if (message.Recver != null && message.hasOwnProperty("Recver"))
                if (typeof message.Recver === "number")
                    object.Recver = options.longs === String ? String(message.Recver) : message.Recver;
                else
                    object.Recver = options.longs === String ? $util.Long.prototype.toString.call(message.Recver) : options.longs === Number ? new $util.LongBits(message.Recver.low >>> 0, message.Recver.high >>> 0).toNumber() : message.Recver;
            if (message.RecverName != null && message.hasOwnProperty("RecverName"))
                object.RecverName = message.RecverName;
            if (message.MessageType != null && message.hasOwnProperty("MessageType"))
                object.MessageType = message.MessageType;
            if (message.Message != null && message.hasOwnProperty("Message"))
                object.Message = message.Message;
            return object;
        };

        /**
         * Converts this W_C_ChatMessage to JSON.
         * @function toJSON
         * @memberof message.W_C_ChatMessage
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        W_C_ChatMessage.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return W_C_ChatMessage;
    })();

    message.W_C_Test = (function() {

        /**
         * Properties of a W_C_Test.
         * @memberof message
         * @interface IW_C_Test
         * @property {Array.<number>|null} [Recv] W_C_Test Recv
         */

        /**
         * Constructs a new W_C_Test.
         * @memberof message
         * @classdesc Represents a W_C_Test.
         * @implements IW_C_Test
         * @constructor
         * @param {message.IW_C_Test=} [properties] Properties to set
         */
        function W_C_Test(properties) {
            this.Recv = [];
            if (properties)
                for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * W_C_Test Recv.
         * @member {Array.<number>} Recv
         * @memberof message.W_C_Test
         * @instance
         */
        W_C_Test.prototype.Recv = $util.emptyArray;

        /**
         * Creates a new W_C_Test instance using the specified properties.
         * @function create
         * @memberof message.W_C_Test
         * @static
         * @param {message.IW_C_Test=} [properties] Properties to set
         * @returns {message.W_C_Test} W_C_Test instance
         */
        W_C_Test.create = function create(properties) {
            return new W_C_Test(properties);
        };

        /**
         * Encodes the specified W_C_Test message. Does not implicitly {@link message.W_C_Test.verify|verify} messages.
         * @function encode
         * @memberof message.W_C_Test
         * @static
         * @param {message.IW_C_Test} message W_C_Test message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        W_C_Test.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.Recv != null && message.Recv.length) {
                writer.uint32(/* id 1, wireType 2 =*/10).fork();
                for (var i = 0; i < message.Recv.length; ++i)
                    writer.int32(message.Recv[i]);
                writer.ldelim();
            }
            return writer;
        };

        /**
         * Encodes the specified W_C_Test message, length delimited. Does not implicitly {@link message.W_C_Test.verify|verify} messages.
         * @function encodeDelimited
         * @memberof message.W_C_Test
         * @static
         * @param {message.IW_C_Test} message W_C_Test message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        W_C_Test.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a W_C_Test message from the specified reader or buffer.
         * @function decode
         * @memberof message.W_C_Test
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {message.W_C_Test} W_C_Test
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        W_C_Test.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            var end = length === undefined ? reader.len : reader.pos + length, message = new $root.message.W_C_Test();
            while (reader.pos < end) {
                var tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    if (!(message.Recv && message.Recv.length))
                        message.Recv = [];
                    if ((tag & 7) === 2) {
                        var end2 = reader.uint32() + reader.pos;
                        while (reader.pos < end2)
                            message.Recv.push(reader.int32());
                    } else
                        message.Recv.push(reader.int32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a W_C_Test message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof message.W_C_Test
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {message.W_C_Test} W_C_Test
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        W_C_Test.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a W_C_Test message.
         * @function verify
         * @memberof message.W_C_Test
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        W_C_Test.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.Recv != null && message.hasOwnProperty("Recv")) {
                if (!Array.isArray(message.Recv))
                    return "Recv: array expected";
                for (var i = 0; i < message.Recv.length; ++i)
                    if (!$util.isInteger(message.Recv[i]))
                        return "Recv: integer[] expected";
            }
            return null;
        };

        /**
         * Creates a W_C_Test message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof message.W_C_Test
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {message.W_C_Test} W_C_Test
         */
        W_C_Test.fromObject = function fromObject(object) {
            if (object instanceof $root.message.W_C_Test)
                return object;
            var message = new $root.message.W_C_Test();
            if (object.Recv) {
                if (!Array.isArray(object.Recv))
                    throw TypeError(".message.W_C_Test.Recv: array expected");
                message.Recv = [];
                for (var i = 0; i < object.Recv.length; ++i)
                    message.Recv[i] = object.Recv[i] | 0;
            }
            return message;
        };

        /**
         * Creates a plain object from a W_C_Test message. Also converts values to other types if specified.
         * @function toObject
         * @memberof message.W_C_Test
         * @static
         * @param {message.W_C_Test} message W_C_Test
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        W_C_Test.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            var object = {};
            if (options.arrays || options.defaults)
                object.Recv = [];
            if (message.Recv && message.Recv.length) {
                object.Recv = [];
                for (var j = 0; j < message.Recv.length; ++j)
                    object.Recv[j] = message.Recv[j];
            }
            return object;
        };

        /**
         * Converts this W_C_Test to JSON.
         * @function toJSON
         * @memberof message.W_C_Test
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        W_C_Test.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return W_C_Test;
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
