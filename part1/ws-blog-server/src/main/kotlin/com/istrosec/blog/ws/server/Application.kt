package com.istrosec.blog.ws.server

import io.ktor.application.*
import io.ktor.features.*
import io.ktor.http.*
import io.ktor.response.*
import io.ktor.routing.*
import io.ktor.server.engine.*
import io.ktor.server.netty.*
import io.ktor.websocket.*

fun main() {
    embeddedServer(Netty, port = 8080, host = "0.0.0.0") {
        installJson()
        install(WebSockets)
        install(CORS) {
            method(HttpMethod.Options)
            method(HttpMethod.Put)
            method(HttpMethod.Delete)
            method(HttpMethod.Patch)
            header(HttpHeaders.Authorization)
            allowCredentials = true
            anyHost() // @TODO: Don't do this in production if possible. Try to limit it.
        }
        routing {
            get("/") {
                call.respondText("Hello World")
            }
            agent()
            client()
        }
    }.start(wait = true)
}

private fun Routing.agent() {
    webSocket("/api/agent/shell") { // websocketSession
        passShellMessages(this, agents, clients)
    }
}

private fun Routing.client() {
    webSocket("/api/client/shell") { // websocketSession
        passShellMessages(this, clients, agents)
    }
    get("/api/agents") {
        call.respond(agents.list())
    }
}