package io.tetrate.log4shell.vulnerable;

import org.apache.logging.log4j.Level;
import org.apache.logging.log4j.core.config.Configurator;
import org.eclipse.jetty.server.Server;
import org.eclipse.jetty.servlet.ServletContextHandler;

public final class App {

    public static void main(String[] args) throws Exception {
        Configurator.setRootLevel(Level.INFO);
        Configurator.setLevel(GreetingsServlet.class.getName(), Level.DEBUG);

        Server server = new Server(8080);
        ServletContextHandler handler = new ServletContextHandler(server, "/");
        handler.addServlet(GreetingsServlet.class, "/");
        server.start();
    }

}
