FROM nginx:stable-alpine

ENV PAGE_TITLE=odysseia
RUN rm -rf /etc/nginx/conf.d
COPY conf /etc/nginx

COPY *.html /usr/share/nginx/html/

COPY yaml/*.yaml /usr/share/nginx/html/swagger/

EXPOSE 8080
CMD ["nginx", "-g", "daemon off;"]
