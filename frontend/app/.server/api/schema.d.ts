/**
 * This file was auto-generated by openapi-typescript.
 * Do not make direct changes to the file.
 */

export interface paths {
    "/login": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** User login */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody: {
                content: {
                    "application/json": {
                        /** @example john */
                        username: string;
                        /** @example password123 */
                        password: string;
                    };
                };
            };
            responses: {
                /** @description Successfully authenticated */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            /** @example xxxxx.xxxxx.xxxxx */
                            token: string;
                        };
                    };
                };
                /** @description Invalid username or password */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            /** @example Invalid credentials */
                            message: string;
                        };
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/token": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get a short-lived access token */
        get: {
            parameters: {
                query?: never;
                header: {
                    Authorization: string;
                };
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Successfully authenticated */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            /** @example xxxxx.xxxxx.xxxxx */
                            token: string;
                        };
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            /** @example Unauthorized */
                            message: string;
                        };
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/games": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** List games */
        get: {
            parameters: {
                query?: {
                    player_id?: number;
                };
                header: {
                    Authorization: string;
                };
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description List of games */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            games: components["schemas"]["Game"][];
                        };
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            /** @example Unauthorized */
                            message: string;
                        };
                    };
                };
                /** @description Forbidden */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            /** @example Forbidden operation */
                            message: string;
                        };
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    [path: `/games/${integer}`]: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get a game */
        get: {
            parameters: {
                query?: never;
                header: {
                    Authorization: string;
                };
                path: {
                    game_id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description A game */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            game: components["schemas"]["Game"];
                        };
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            /** @example Unauthorized */
                            message: string;
                        };
                    };
                };
                /** @description Forbidden */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            /** @example Forbidden operation */
                            message: string;
                        };
                    };
                };
                /** @description Not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            /** @example Not found */
                            message: string;
                        };
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/admin/users": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** List all users */
        get: {
            parameters: {
                query?: never;
                header: {
                    Authorization: string;
                };
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description List of users */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            users: components["schemas"]["User"][];
                        };
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            /** @example Unauthorized */
                            message: string;
                        };
                    };
                };
                /** @description Forbidden */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            /** @example Forbidden operation */
                            message: string;
                        };
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/admin/games": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** List games */
        get: {
            parameters: {
                query?: never;
                header: {
                    Authorization: string;
                };
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description List of games */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            games: components["schemas"]["Game"][];
                        };
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            /** @example Unauthorized */
                            message: string;
                        };
                    };
                };
                /** @description Forbidden */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            /** @example Forbidden operation */
                            message: string;
                        };
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    [path: `/admin/games/${integer}`]: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get a game */
        get: {
            parameters: {
                query?: never;
                header: {
                    Authorization: string;
                };
                path: {
                    game_id: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description A game */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            game: components["schemas"]["Game"];
                        };
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            /** @example Unauthorized */
                            message: string;
                        };
                    };
                };
                /** @description Forbidden */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            /** @example Forbidden operation */
                            message: string;
                        };
                    };
                };
                /** @description Not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            /** @example Not found */
                            message: string;
                        };
                    };
                };
            };
        };
        /** Update a game */
        put: {
            parameters: {
                query?: never;
                header: {
                    Authorization: string;
                };
                path: {
                    game_id: number;
                };
                cookie?: never;
            };
            requestBody: {
                content: {
                    "application/json": {
                        /**
                         * @example closed
                         * @enum {string}
                         */
                        state?: "closed" | "waiting_entries" | "waiting_start" | "prepare" | "starting" | "gaming" | "finished";
                        /** @example Game 1 */
                        display_name?: string;
                        /** @example 360 */
                        duration_seconds?: number;
                        /** @example 946684800 */
                        started_at?: number | null;
                        /** @example 1 */
                        problem_id?: number | null;
                    };
                };
            };
            responses: {
                /** @description Successfully updated */
                204: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content?: never;
                };
                /** @description Invalid request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            /** @example Invalid request */
                            message: string;
                        };
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            /** @example Unauthorized */
                            message: string;
                        };
                    };
                };
                /** @description Forbidden */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            /** @example Forbidden operation */
                            message: string;
                        };
                    };
                };
                /** @description Not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": {
                            /** @example Not found */
                            message: string;
                        };
                    };
                };
            };
        };
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
}
export type webhooks = Record<string, never>;
export interface components {
    schemas: {
        User: {
            /** @example 123 */
            user_id: number;
            /** @example john */
            username: string;
            /** @example John Doe */
            display_name: string;
            /** @example /images/john.jpg */
            icon_path?: string;
            /** @example false */
            is_admin: boolean;
        };
        Game: {
            /** @example 1 */
            game_id: number;
            /**
             * @example closed
             * @enum {string}
             */
            state: "closed" | "waiting_entries" | "waiting_start" | "prepare" | "starting" | "gaming" | "finished";
            /** @example Game 1 */
            display_name: string;
            /** @example 360 */
            duration_seconds: number;
            /** @example 946684800 */
            started_at?: number;
            problem?: components["schemas"]["Problem"];
        };
        Problem: {
            /** @example 1 */
            problem_id: number;
            /** @example Problem 1 */
            title: string;
            /** @example This is a problem */
            description: string;
        };
        GamePlayerMessage: components["schemas"]["GamePlayerMessageS2C"] | components["schemas"]["GamePlayerMessageC2S"];
        GamePlayerMessageS2C: components["schemas"]["GamePlayerMessageS2CPrepare"] | components["schemas"]["GamePlayerMessageS2CStart"] | components["schemas"]["GamePlayerMessageS2CExecResult"];
        GamePlayerMessageS2CPrepare: {
            /** @constant */
            type: "player:s2c:prepare";
            data: components["schemas"]["GamePlayerMessageS2CPreparePayload"];
        };
        GamePlayerMessageS2CPreparePayload: {
            problem: components["schemas"]["Problem"];
        };
        GamePlayerMessageS2CStart: {
            /** @constant */
            type: "player:s2c:start";
            data: components["schemas"]["GamePlayerMessageS2CStartPayload"];
        };
        GamePlayerMessageS2CStartPayload: {
            /** @example 946684800 */
            start_at: number;
        };
        GamePlayerMessageS2CExecResult: {
            /** @constant */
            type: "player:s2c:execresult";
            data: components["schemas"]["GamePlayerMessageS2CExecResultPayload"];
        };
        GamePlayerMessageS2CExecResultPayload: {
            /**
             * @example success
             * @enum {string}
             */
            status: "success";
            /** @example 100 */
            score: number | null;
        };
        GamePlayerMessageC2S: components["schemas"]["GamePlayerMessageC2SEntry"] | components["schemas"]["GamePlayerMessageC2SReady"] | components["schemas"]["GamePlayerMessageC2SCode"];
        GamePlayerMessageC2SEntry: {
            /** @constant */
            type: "player:c2s:entry";
        };
        GamePlayerMessageC2SReady: {
            /** @constant */
            type: "player:c2s:ready";
        };
        GamePlayerMessageC2SCode: {
            /** @constant */
            type: "player:c2s:code";
            data: components["schemas"]["GamePlayerMessageC2SCodePayload"];
        };
        GamePlayerMessageC2SCodePayload: {
            /** @example print('Hello, world!') */
            code: string;
        };
        GameWatcherMessage: components["schemas"]["GameWatcherMessageS2C"];
        GameWatcherMessageS2C: components["schemas"]["GameWatcherMessageS2CStart"] | components["schemas"]["GameWatcherMessageS2CCode"] | components["schemas"]["GameWatcherMessageS2CExecResult"];
        GameWatcherMessageS2CStart: {
            /** @constant */
            type: "watcher:s2c:start";
            data: components["schemas"]["GameWatcherMessageS2CStartPayload"];
        };
        GameWatcherMessageS2CStartPayload: {
            /** @example 946684800 */
            start_at: number;
        };
        GameWatcherMessageS2CCode: {
            /** @constant */
            type: "watcher:s2c:code";
            data: components["schemas"]["GameWatcherMessageS2CCodePayload"];
        };
        GameWatcherMessageS2CCodePayload: {
            /** @example 1 */
            player_id: number;
            /** @example print('Hello, world!') */
            code: string;
        };
        GameWatcherMessageS2CExecResult: {
            /** @constant */
            type: "watcher:s2c:execresult";
            data: components["schemas"]["GameWatcherMessageS2CExecResultPayload"];
        };
        GameWatcherMessageS2CExecResultPayload: {
            /** @example 1 */
            player_id: number;
            /**
             * @example success
             * @enum {string}
             */
            status: "success";
            /** @example 100 */
            score: number | null;
            /** @example Hello, world! */
            stdout: string;
            /** @example  */
            stderr: string;
        };
    };
    responses: never;
    parameters: never;
    requestBodies: never;
    headers: never;
    pathItems: never;
}
export type $defs = Record<string, never>;
export type operations = Record<string, never>;
