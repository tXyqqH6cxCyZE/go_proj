// Autogenerated by Thrift Compiler (0.9.2)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package hello

import (
	"bytes"
	"fmt"
	"micode.be.xiaomi.com/systech/soa/filter"
	"micode.be.xiaomi.com/systech/soa/thrift"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = filter.BEFORE
var _ = bytes.Equal

type HelloService interface {
	// Parameters:
	//  - Name
	HelloWorld(ctx *thrift.XContext, name string) (r *Result_, errRet error)
}

type HelloServiceClient struct {
	Transport       thrift.TTransport
	ProtocolFactory thrift.TProtocolFactory
	InputProtocol   thrift.TProtocol
	OutputProtocol  thrift.TProtocol
	SeqId           int32
}

func NewHelloServiceClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *HelloServiceClient {
	return &HelloServiceClient{Transport: t,
		ProtocolFactory: f,
		InputProtocol:   f.GetProtocol(t),
		OutputProtocol:  f.GetProtocol(t),
		SeqId:           0,
	}
}

func NewHelloServiceClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *HelloServiceClient {
	return &HelloServiceClient{Transport: t,
		ProtocolFactory: nil,
		InputProtocol:   iprot,
		OutputProtocol:  oprot,
		SeqId:           0,
	}
}

// Parameters:
//  - Name
func (p *HelloServiceClient) HelloWorld(ctx *thrift.XContext, name string) (r *Result_, errRet error) {
	if errRet = p.sendHelloWorld(ctx, name); errRet != nil {
		return
	}
	return p.recvHelloWorld()
}

func (p *HelloServiceClient) sendHelloWorld(ctx *thrift.XContext, name string) (err error) {
	oprot := p.OutputProtocol
	xmheaderTrans, ok := p.Transport.(*thrift.XmHeaderTransport)
	if ok && ctx != nil {
		xmheaderTrans.SetAppId(ctx.GetAppId())
		xmheaderTrans.SetLogId(ctx.GetLogId())
		xmheaderTrans.SetRpcId(ctx.GetRpcId())
	}

	if oprot == nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	if err = oprot.WriteMessageBegin("HelloWorld", thrift.CALL, p.SeqId); err != nil {
		return
	}
	args := HelloWorldArgs{
		Name: name,
	}
	if err = args.Write(oprot); err != nil {
		return
	}
	if err = oprot.WriteMessageEnd(); err != nil {
		return
	}
	return oprot.Flush()
}

func (p *HelloServiceClient) recvHelloWorld() (value *Result_, err error) {
	iprot := p.InputProtocol
	if iprot == nil {
		iprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.InputProtocol = iprot
	}
	_, mTypeId, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return
	}
	if mTypeId == thrift.EXCEPTION {
		error0 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
		var error1 error
		error1, err = error0.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error1
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "HelloWorld failed: out of sequence response")
		return
	}
	result := HelloWorldResult{}
	if err = result.Read(iprot); err != nil {
		return
	}
	if err = iprot.ReadMessageEnd(); err != nil {
		return
	}
	value = result.GetSuccess()
	return
}

type HelloServiceProcessor struct {
	processorMap map[string]thrift.TProcessorFunction
	filters      map[string]map[filter.FilterType][]filter.Filter
	handler      HelloService
}

func (p *HelloServiceProcessor) AddToProcessorMap(key string, processor thrift.TProcessorFunction) {
	p.processorMap[key] = processor
}

func (p *HelloServiceProcessor) GetProcessorFunction(key string) (processor thrift.TProcessorFunction, ok bool) {
	processor, ok = p.processorMap[key]
	return processor, ok
}

func (p *HelloServiceProcessor) ProcessorMap() map[string]thrift.TProcessorFunction {
	return p.processorMap
}

func NewHelloServiceProcessor(handler HelloService) *HelloServiceProcessor {

	self2 := &HelloServiceProcessor{handler: handler, processorMap: make(map[string]thrift.TProcessorFunction)}
	self2.processorMap["HelloWorld"] = &helloServiceProcessorHelloWorld{handler: handler, processor: self2}
	return self2
}

func (p *HelloServiceProcessor) Process(iprot, oprot thrift.TProtocol, ctx *thrift.XContext) (success bool, err thrift.TException) {
	name, _, seqId, err := iprot.ReadMessageBegin()
	ctx.StartTime("total")
	if err != nil {
		return false, err
	}
	if processor, ok := p.GetProcessorFunction(name); ok {
		ctx.SetMethod(name)
		return processor.Process(seqId, iprot, oprot, ctx)
	}
	iprot.Skip(thrift.STRUCT)
	iprot.ReadMessageEnd()
	x3 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "Unknown function "+name)
	oprot.WriteMessageBegin(name, thrift.EXCEPTION, seqId)
	x3.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush()
	return false, x3

}

type helloServiceProcessorHelloWorld struct {
	handler   HelloService
	processor *HelloServiceProcessor
}

func (p *helloServiceProcessorHelloWorld) Process(seqId int32, iprot, oprot thrift.TProtocol, ctx *thrift.XContext) (success bool, err thrift.TException) {
	args := HelloWorldArgs{}
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("HelloWorld", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return false, err
	}

	iprot.ReadMessageEnd()

	trans := iprot.Transport()
	xmheaderTrans, ok := trans.(*thrift.XmHeaderTransport)
	if ok {
		header := xmheaderTrans.GetXmHeader()
		ctx.SetParentRpcId(header.RpcId)
		ctx.SetLogId(header.LogId)
		ctx.SetAppId(header.AppId)
		err = ctx.BeforeProcess()
		if err != nil {
			return
		}
	}

	result := HelloWorldResult{}
	var retval *Result_ = &Result_{}
	var err2 error
	beforeFilters := p.processor.getBeforeFilters("HelloWorld")
	ret := false
	for _, f := range beforeFilters {
		if f.Filter(ctx, "HelloWorld", []interface{}{&args.Name}, &filter.Response{Val: retval}) {
			ret = true
			break
		}
	}
	if !ret {
		if retval, err2 = p.handler.HelloWorld(ctx, args.Name); err2 != nil {
			afterFilters := p.processor.getAfterFilters("HelloWorld")
			for _, f := range afterFilters {
				if f.Filter(ctx, "HelloWorld", []interface{}{&args.Name}, &filter.Response{Val: retval}) {
					break
				}
			}
			x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing HelloWorld: "+err2.Error())
			oprot.WriteMessageBegin("HelloWorld", thrift.EXCEPTION, seqId)
			x.Write(oprot)
			oprot.WriteMessageEnd()
			oprot.Flush()
			return true, x
		} else {
			afterFilters := p.processor.getAfterFilters("HelloWorld")
			for _, f := range afterFilters {
				if f.Filter(ctx, "HelloWorld", []interface{}{&args.Name}, &filter.Response{Val: retval}) {
					break
				}
			}
			result.Success = retval
		}
	} else {
		result.Success = retval
	}

	if thrift.IsPrintInput("HelloWorld") == true {
		ctx.Set("input", args)
	}
	if thrift.IsPrintOutput("HelloWorld") == true {
		ctx.Set("output", result)
	}

	if err2 = oprot.WriteMessageBegin("HelloWorld", thrift.REPLY, seqId); err2 != nil {
		err = err2
	}
	if err2 = result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.Flush(); err == nil && err2 != nil {
		err = err2
	}
	if err != nil {
		return
	}
	return true, err
}

func (p *HelloServiceProcessor) registerFilter(handler string, typ filter.FilterType, f filter.Filter) {
	if p.filters == nil {
		p.filters = make(map[string]map[filter.FilterType][]filter.Filter)
	}
	if _, ok1 := p.filters[handler]; !ok1 {
		p.filters[handler] = make(map[filter.FilterType][]filter.Filter)
	}
	if _, ok2 := p.filters[handler][typ]; !ok2 {
		p.filters[handler][typ] = make([]filter.Filter, 0)
	}
	p.filters[handler][typ] = append(p.filters[handler][typ], f)
}
func (p *HelloServiceProcessor) RegisterFilter(handlers []string, typ filter.FilterType, f filter.Filter) {
	for _, handler := range handlers {
		p.registerFilter(handler, typ, f)
	}
}
func (p *HelloServiceProcessor) getBeforeFilters(handler string) []filter.Filter {
	return p.getFilters(handler, filter.BEFORE)
}
func (p *HelloServiceProcessor) getAfterFilters(handler string) []filter.Filter {
	return p.getFilters(handler, filter.AFTER)
}
func (p *HelloServiceProcessor) getFilters(handler string, typ filter.FilterType) []filter.Filter {
	if _, ok1 := p.filters[handler]; !ok1 {
		return []filter.Filter{}
	}
	if _, ok2 := p.filters[handler][typ]; !ok2 {
		return []filter.Filter{}
	}
	return p.filters[handler][typ]
}

// HELPER FUNCTIONS AND STRUCTURES

type HelloWorldArgs struct {
	Name string `thrift:"name,1" json:"name"`
}

func NewHelloWorldArgs() *HelloWorldArgs {
	return &HelloWorldArgs{}
}

func (p *HelloWorldArgs) GetName() string {
	return p.Name
}
func (p *HelloWorldArgs) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error: %s", p, err)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *HelloWorldArgs) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 1: %s", err)
	} else {
		p.Name = v
	}
	return nil
}

func (p *HelloWorldArgs) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("HelloWorld_args"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("write struct stop error: %s", err)
	}
	return nil
}

func (p *HelloWorldArgs) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("name", thrift.STRING, 1); err != nil {
		return fmt.Errorf("%T write field begin error 1:name: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Name)); err != nil {
		return fmt.Errorf("%T.name (1) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:name: %s", p, err)
	}
	return err
}

func (p *HelloWorldArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("HelloWorldArgs(%+v)", *p)
}

type HelloWorldResult struct {
	Success *Result_ `thrift:"success,0" json:"success"`
}

func NewHelloWorldResult() *HelloWorldResult {
	return &HelloWorldResult{}
}

var HelloWorldResult_Success_DEFAULT *Result_

func (p *HelloWorldResult) GetSuccess() *Result_ {
	if !p.IsSetSuccess() {
		return HelloWorldResult_Success_DEFAULT
	}
	return p.Success
}
func (p *HelloWorldResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *HelloWorldResult) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error: %s", p, err)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 0:
			if err := p.ReadField0(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *HelloWorldResult) ReadField0(iprot thrift.TProtocol) error {
	p.Success = &Result_{}
	if err := p.Success.Read(iprot); err != nil {
		return fmt.Errorf("%T error reading struct: %s", p.Success, err)
	}
	return nil
}

func (p *HelloWorldResult) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("HelloWorld_result"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField0(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("write struct stop error: %s", err)
	}
	return nil
}

func (p *HelloWorldResult) writeField0(oprot thrift.TProtocol) (err error) {
	if p.IsSetSuccess() {
		if err := oprot.WriteFieldBegin("success", thrift.STRUCT, 0); err != nil {
			return fmt.Errorf("%T write field begin error 0:success: %s", p, err)
		}
		if err := p.Success.Write(oprot); err != nil {
			return fmt.Errorf("%T error writing struct: %s", p.Success, err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return fmt.Errorf("%T write field end error 0:success: %s", p, err)
		}
	}
	return err
}

func (p *HelloWorldResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("HelloWorldResult(%+v)", *p)
}
