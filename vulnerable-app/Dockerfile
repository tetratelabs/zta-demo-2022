FROM maven:3-jdk-8 as builder
RUN mkdir /build/
WORKDIR /build
COPY ./ ./
RUN mvn clean compile assembly:single

FROM openjdk:8u102-jdk
COPY --from=builder /build/target/vulnerable-1.0-SNAPSHOT-jar-with-dependencies.jar .
EXPOSE 8080
#ENV LOG4J_FORMAT_MSG_NO_LOOKUPS=false
#CMD java -Dcom.sun.jndi.ldap.object.trustURLCodebase=true -Dcom.sun.jndi.rmi.object.trustURLCodebase=true -jar vulnerable-1.0-SNAPSHOT-jar-with-dependencies.jar
CMD java -jar vulnerable-1.0-SNAPSHOT-jar-with-dependencies.jar
