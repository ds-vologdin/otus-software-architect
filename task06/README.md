# Architecture

Below we will describe the communication and IDL schemes for different architecture variants. We will consider schemes for only creating an order.

## RESTful

## IDL schemes

- [Order service](https://github.com/ds-vologdin/otus-software-architect/blob/main/task06/arch/restful/order.yml)
- [Billing service](https://github.com/ds-vologdin/otus-software-architect/blob/main/task06/arch/restful/billing.yml)
- [Notification service](https://github.com/ds-vologdin/otus-software-architect/blob/main/task06/arch/restful/notification.yml)
- [Account service](https://github.com/ds-vologdin/otus-software-architect/blob/main/task06/arch/restful/account.yml)

## Communication scheme
![restful](https://github.com/ds-vologdin/otus-software-architect/blob/main/task06/images/restful.jpg)

## Notes

Notification service should use a cache to reduce the count of requests to Account service.

### Event Notifications

![event notification](https://github.com/ds-vologdin/otus-software-architect/blob/main/task06/images/event_notification.jpg)

## Event Collaboration

![event collaboration](https://github.com/ds-vologdin/otus-software-architect/blob/main/task06/images/event_collaboration.jpg)
