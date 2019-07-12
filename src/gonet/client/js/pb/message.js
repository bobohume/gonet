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
