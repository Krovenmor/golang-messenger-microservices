FROM nginx:1.27-alpine AS runtime

COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 443
EXPOSE 80

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://127.0.0.1/health/check >/dev/null 2>&1 || exit 1

CMD ["nginx", "-g", "daemon off;"]