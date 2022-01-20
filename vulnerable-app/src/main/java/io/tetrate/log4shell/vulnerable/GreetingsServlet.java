package io.tetrate.log4shell.vulnerable;

import java.io.IOException;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import com.nimbusds.jose.JWSObject;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class GreetingsServlet extends HttpServlet {

    private static final long serialVersionUID = 1L;

    private static Logger log = LogManager.getLogger(GreetingsServlet.class.getName());

    @Override
    protected void doGet(HttpServletRequest request, HttpServletResponse response)
            throws ServletException, IOException {

        String token = getToken(request);
        String subject = getSubject(token);

        // This is the  exploitable log line
        log.info("welcoming user: " + subject);

        response.setContentType("text/plain");
        response.setStatus(HttpServletResponse.SC_OK);
        response.getWriter().println("Welcome, " + subject + "!");
    }

    private String getToken(HttpServletRequest req) {
        String auth = req.getHeader("Authorization");
        if (auth == null) {
            log.debug("no authorization header present");
            return null;
        }

        String[] parts = auth.split(" ");
        if (parts.length != 2) {
            log.debug("invalid authorization header value");
            return null;
        }

        return parts[1];
    }

    private String getSubject(String token) {
        String subject = "anonymous";

        if (token != null) {
            try {
                JWSObject jet = JWSObject.parse(token);
                subject = (String) jet.getPayload().toJSONObject().get("sub");
            } catch (Exception ex) {
                ex.printStackTrace();
                subject = "unknown";
            }
        }

        return subject;
    }
}
