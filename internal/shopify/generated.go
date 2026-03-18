package shopify

import (
	"context"

	"github.com/Khan/genqlient/graphql"
)

type WebhookSubscriptionFormat string

const (
	WebhookSubscriptionFormatJson WebhookSubscriptionFormat = "JSON"
	WebhookSubscriptionFormatXml  WebhookSubscriptionFormat = "XML"
)

type WebhookSubscriptionInput struct {
	Format *WebhookSubscriptionFormat `json:"format"`
	Uri    *string                    `json:"uri"`
}

func (v *WebhookSubscriptionInput) GetFormat() *WebhookSubscriptionFormat { return v.Format }
func (v *WebhookSubscriptionInput) GetUri() *string                       { return v.Uri }

type WebhookSubscriptionTopic string

const (
	WebhookSubscriptionTopicAppUninstalled                        WebhookSubscriptionTopic = "APP_UNINSTALLED"
	WebhookSubscriptionTopicCustomersDataRequest                  WebhookSubscriptionTopic = "CUSTOMERS_DATA_REQUEST"
	WebhookSubscriptionTopicCustomersRedact                       WebhookSubscriptionTopic = "CUSTOMERS_REDACT"
	WebhookSubscriptionTopicShopRedact                            WebhookSubscriptionTopic = "SHOP_REDACT"
	WebhookSubscriptionTopicSubscriptionContractsCreate           WebhookSubscriptionTopic = "SUBSCRIPTION_CONTRACTS_CREATE"
	WebhookSubscriptionTopicSubscriptionContractsUpdate           WebhookSubscriptionTopic = "SUBSCRIPTION_CONTRACTS_UPDATE"
	WebhookSubscriptionTopicSubscriptionContractsActivate         WebhookSubscriptionTopic = "SUBSCRIPTION_CONTRACTS_ACTIVATE"
	WebhookSubscriptionTopicSubscriptionContractsPause            WebhookSubscriptionTopic = "SUBSCRIPTION_CONTRACTS_PAUSE"
	WebhookSubscriptionTopicSubscriptionContractsCancel           WebhookSubscriptionTopic = "SUBSCRIPTION_CONTRACTS_CANCEL"
	WebhookSubscriptionTopicSubscriptionContractsFail             WebhookSubscriptionTopic = "SUBSCRIPTION_CONTRACTS_FAIL"
	WebhookSubscriptionTopicSubscriptionContractsExpire           WebhookSubscriptionTopic = "SUBSCRIPTION_CONTRACTS_EXPIRE"
	WebhookSubscriptionTopicSubscriptionBillingAttemptsSuccess    WebhookSubscriptionTopic = "SUBSCRIPTION_BILLING_ATTEMPTS_SUCCESS"
	WebhookSubscriptionTopicSubscriptionBillingAttemptsFailure    WebhookSubscriptionTopic = "SUBSCRIPTION_BILLING_ATTEMPTS_FAILURE"
	WebhookSubscriptionTopicSubscriptionBillingAttemptsChallenged WebhookSubscriptionTopic = "SUBSCRIPTION_BILLING_ATTEMPTS_CHALLENGED"
	WebhookSubscriptionTopicCustomerPaymentMethodsCreate          WebhookSubscriptionTopic = "CUSTOMER_PAYMENT_METHODS_CREATE"
	WebhookSubscriptionTopicCustomerPaymentMethodsUpdate          WebhookSubscriptionTopic = "CUSTOMER_PAYMENT_METHODS_UPDATE"
	WebhookSubscriptionTopicCustomerPaymentMethodsRevoke          WebhookSubscriptionTopic = "CUSTOMER_PAYMENT_METHODS_REVOKE"
)

type WebhookSubscriptionsByTopicResponse struct {
	WebhookSubscriptions WebhookSubscriptionsByTopicWebhookSubscriptionsWebhookSubscriptionConnection `json:"webhookSubscriptions"`
}

type WebhookSubscriptionsByTopicWebhookSubscriptionsWebhookSubscriptionConnection struct {
	Nodes []WebhookSubscriptionsByTopicWebhookSubscriptionsWebhookSubscriptionConnectionNodesWebhookSubscription `json:"nodes"`
}

type WebhookSubscriptionsByTopicWebhookSubscriptionsWebhookSubscriptionConnectionNodesWebhookSubscription struct {
	Id    string                   `json:"id"`
	Topic WebhookSubscriptionTopic `json:"topic"`
	Uri   string                   `json:"uri"`
}

func (v *WebhookSubscriptionsByTopicWebhookSubscriptionsWebhookSubscriptionConnectionNodesWebhookSubscription) GetId() string {
	return v.Id
}

func (v *WebhookSubscriptionsByTopicWebhookSubscriptionsWebhookSubscriptionConnectionNodesWebhookSubscription) GetTopic() WebhookSubscriptionTopic {
	return v.Topic
}

func (v *WebhookSubscriptionsByTopicWebhookSubscriptionsWebhookSubscriptionConnectionNodesWebhookSubscription) GetUri() string {
	return v.Uri
}

type WebhookSubscriptionCreateResponse struct {
	WebhookSubscriptionCreate *WebhookSubscriptionCreateWebhookSubscriptionCreateWebhookSubscriptionCreatePayload `json:"webhookSubscriptionCreate"`
}

type WebhookSubscriptionCreateWebhookSubscriptionCreateWebhookSubscriptionCreatePayload struct {
	WebhookSubscription *WebhookSubscriptionCreateWebhookSubscriptionCreateWebhookSubscriptionCreatePayloadWebhookSubscription  `json:"webhookSubscription"`
	UserErrors          []WebhookSubscriptionCreateWebhookSubscriptionCreateWebhookSubscriptionCreatePayloadUserErrorsUserError `json:"userErrors"`
}

type WebhookSubscriptionCreateWebhookSubscriptionCreateWebhookSubscriptionCreatePayloadWebhookSubscription struct {
	Id    string                   `json:"id"`
	Topic WebhookSubscriptionTopic `json:"topic"`
	Uri   string                   `json:"uri"`
}

type WebhookSubscriptionCreateWebhookSubscriptionCreateWebhookSubscriptionCreatePayloadUserErrorsUserError struct {
	Field   []string `json:"field"`
	Message string   `json:"message"`
}

func (v *WebhookSubscriptionCreateWebhookSubscriptionCreateWebhookSubscriptionCreatePayloadUserErrorsUserError) GetField() []string {
	return v.Field
}

func (v *WebhookSubscriptionCreateWebhookSubscriptionCreateWebhookSubscriptionCreatePayloadUserErrorsUserError) GetMessage() string {
	return v.Message
}

type WebhookSubscriptionDeleteResponse struct {
	WebhookSubscriptionDelete *WebhookSubscriptionDeleteWebhookSubscriptionDeleteWebhookSubscriptionDeletePayload `json:"webhookSubscriptionDelete"`
}

type WebhookSubscriptionDeleteWebhookSubscriptionDeleteWebhookSubscriptionDeletePayload struct {
	DeletedWebhookSubscriptionId *string                                                                                                 `json:"deletedWebhookSubscriptionId"`
	UserErrors                   []WebhookSubscriptionDeleteWebhookSubscriptionDeleteWebhookSubscriptionDeletePayloadUserErrorsUserError `json:"userErrors"`
}

type WebhookSubscriptionDeleteWebhookSubscriptionDeleteWebhookSubscriptionDeletePayloadUserErrorsUserError struct {
	Field   []string `json:"field"`
	Message string   `json:"message"`
}

func (v *WebhookSubscriptionDeleteWebhookSubscriptionDeleteWebhookSubscriptionDeletePayloadUserErrorsUserError) GetField() []string {
	return v.Field
}

func (v *WebhookSubscriptionDeleteWebhookSubscriptionDeleteWebhookSubscriptionDeletePayloadUserErrorsUserError) GetMessage() string {
	return v.Message
}

const WebhookSubscriptionCreate_Operation = `
mutation WebhookSubscriptionCreate ($topic: WebhookSubscriptionTopic!, $webhookSubscription: WebhookSubscriptionInput!) {
	webhookSubscriptionCreate(topic: $topic, webhookSubscription: $webhookSubscription) {
		webhookSubscription {
			... WebhookSubscriptionBasic
		}
		userErrors {
			... WebhookUserError
		}
	}
}
fragment WebhookSubscriptionBasic on WebhookSubscription {
	id
	topic
	uri
}
fragment WebhookUserError on UserError {
	field
	message
}
`

const WebhookSubscriptionDelete_Operation = `
mutation WebhookSubscriptionDelete ($id: ID!) {
	webhookSubscriptionDelete(id: $id) {
		deletedWebhookSubscriptionId
		userErrors {
			... WebhookUserError
		}
	}
}
fragment WebhookUserError on UserError {
	field
	message
}
`

const WebhookSubscriptionsByTopic_Operation = `
query WebhookSubscriptionsByTopic ($topics: [WebhookSubscriptionTopic!]) {
	webhookSubscriptions(first: 100, topics: $topics) {
		nodes {
			... WebhookSubscriptionBasic
		}
	}
}
fragment WebhookSubscriptionBasic on WebhookSubscription {
	id
	topic
	uri
}
`

func WebhookSubscriptionCreate(
	ctx_ context.Context,
	client_ graphql.Client,
	topic WebhookSubscriptionTopic,
	webhookSubscription WebhookSubscriptionInput,
) (data_ *WebhookSubscriptionCreateResponse, err_ error) {
	req_ := &graphql.Request{
		OpName: "WebhookSubscriptionCreate",
		Query:  WebhookSubscriptionCreate_Operation,
		Variables: &struct {
			Topic               WebhookSubscriptionTopic `json:"topic"`
			WebhookSubscription WebhookSubscriptionInput `json:"webhookSubscription"`
		}{
			Topic:               topic,
			WebhookSubscription: webhookSubscription,
		},
	}

	data_ = &WebhookSubscriptionCreateResponse{}
	resp_ := &graphql.Response{Data: data_}

	err_ = client_.MakeRequest(ctx_, req_, resp_)
	return data_, err_
}

func WebhookSubscriptionDelete(
	ctx_ context.Context,
	client_ graphql.Client,
	id string,
) (data_ *WebhookSubscriptionDeleteResponse, err_ error) {
	req_ := &graphql.Request{
		OpName: "WebhookSubscriptionDelete",
		Query:  WebhookSubscriptionDelete_Operation,
		Variables: &struct {
			ID string `json:"id"`
		}{
			ID: id,
		},
	}

	data_ = &WebhookSubscriptionDeleteResponse{}
	resp_ := &graphql.Response{Data: data_}

	err_ = client_.MakeRequest(ctx_, req_, resp_)
	return data_, err_
}

func WebhookSubscriptionsByTopic(
	ctx_ context.Context,
	client_ graphql.Client,
	topics []WebhookSubscriptionTopic,
) (data_ *WebhookSubscriptionsByTopicResponse, err_ error) {
	req_ := &graphql.Request{
		OpName: "WebhookSubscriptionsByTopic",
		Query:  WebhookSubscriptionsByTopic_Operation,
		Variables: &struct {
			Topics []WebhookSubscriptionTopic `json:"topics"`
		}{
			Topics: topics,
		},
	}

	data_ = &WebhookSubscriptionsByTopicResponse{}
	resp_ := &graphql.Response{Data: data_}

	err_ = client_.MakeRequest(ctx_, req_, resp_)
	return data_, err_
}
