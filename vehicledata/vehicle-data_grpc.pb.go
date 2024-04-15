// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: vehicledata/vehicle-data.proto

package vehicledata

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	VehicleData_GetVehicleDataByRegistration_FullMethodName  = "/vehicledata.VehicleData/GetVehicleDataByRegistration"
	VehicleData_GetVehicleDataByChassisEngine_FullMethodName = "/vehicledata.VehicleData/GetVehicleDataByChassisEngine"
	VehicleData_GetDriverDataByLicenseNumber_FullMethodName  = "/vehicledata.VehicleData/GetDriverDataByLicenseNumber"
)

// VehicleDataClient is the client API for VehicleData service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VehicleDataClient interface {
	GetVehicleDataByRegistration(ctx context.Context, in *RegistrationRequest, opts ...grpc.CallOption) (*VehicleInfo, error)
	GetVehicleDataByChassisEngine(ctx context.Context, in *ChassisEngineNumberRequest, opts ...grpc.CallOption) (*VehicleInfo, error)
	GetDriverDataByLicenseNumber(ctx context.Context, in *LicenseRequest, opts ...grpc.CallOption) (*DriverInfo, error)
}

type vehicleDataClient struct {
	cc grpc.ClientConnInterface
}

func NewVehicleDataClient(cc grpc.ClientConnInterface) VehicleDataClient {
	return &vehicleDataClient{cc}
}

func (c *vehicleDataClient) GetVehicleDataByRegistration(ctx context.Context, in *RegistrationRequest, opts ...grpc.CallOption) (*VehicleInfo, error) {
	out := new(VehicleInfo)
	err := c.cc.Invoke(ctx, VehicleData_GetVehicleDataByRegistration_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vehicleDataClient) GetVehicleDataByChassisEngine(ctx context.Context, in *ChassisEngineNumberRequest, opts ...grpc.CallOption) (*VehicleInfo, error) {
	out := new(VehicleInfo)
	err := c.cc.Invoke(ctx, VehicleData_GetVehicleDataByChassisEngine_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vehicleDataClient) GetDriverDataByLicenseNumber(ctx context.Context, in *LicenseRequest, opts ...grpc.CallOption) (*DriverInfo, error) {
	out := new(DriverInfo)
	err := c.cc.Invoke(ctx, VehicleData_GetDriverDataByLicenseNumber_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VehicleDataServer is the server API for VehicleData service.
// All implementations must embed UnimplementedVehicleDataServer
// for forward compatibility
type VehicleDataServer interface {
	GetVehicleDataByRegistration(context.Context, *RegistrationRequest) (*VehicleInfo, error)
	GetVehicleDataByChassisEngine(context.Context, *ChassisEngineNumberRequest) (*VehicleInfo, error)
	GetDriverDataByLicenseNumber(context.Context, *LicenseRequest) (*DriverInfo, error)
	mustEmbedUnimplementedVehicleDataServer()
}

// UnimplementedVehicleDataServer must be embedded to have forward compatible implementations.
type UnimplementedVehicleDataServer struct {
}

func (UnimplementedVehicleDataServer) GetVehicleDataByRegistration(context.Context, *RegistrationRequest) (*VehicleInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVehicleDataByRegistration not implemented")
}
func (UnimplementedVehicleDataServer) GetVehicleDataByChassisEngine(context.Context, *ChassisEngineNumberRequest) (*VehicleInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVehicleDataByChassisEngine not implemented")
}
func (UnimplementedVehicleDataServer) GetDriverDataByLicenseNumber(context.Context, *LicenseRequest) (*DriverInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDriverDataByLicenseNumber not implemented")
}
func (UnimplementedVehicleDataServer) mustEmbedUnimplementedVehicleDataServer() {}

// UnsafeVehicleDataServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VehicleDataServer will
// result in compilation errors.
type UnsafeVehicleDataServer interface {
	mustEmbedUnimplementedVehicleDataServer()
}

func RegisterVehicleDataServer(s grpc.ServiceRegistrar, srv VehicleDataServer) {
	s.RegisterService(&VehicleData_ServiceDesc, srv)
}

func _VehicleData_GetVehicleDataByRegistration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegistrationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VehicleDataServer).GetVehicleDataByRegistration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VehicleData_GetVehicleDataByRegistration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VehicleDataServer).GetVehicleDataByRegistration(ctx, req.(*RegistrationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VehicleData_GetVehicleDataByChassisEngine_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChassisEngineNumberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VehicleDataServer).GetVehicleDataByChassisEngine(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VehicleData_GetVehicleDataByChassisEngine_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VehicleDataServer).GetVehicleDataByChassisEngine(ctx, req.(*ChassisEngineNumberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VehicleData_GetDriverDataByLicenseNumber_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LicenseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VehicleDataServer).GetDriverDataByLicenseNumber(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VehicleData_GetDriverDataByLicenseNumber_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VehicleDataServer).GetDriverDataByLicenseNumber(ctx, req.(*LicenseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// VehicleData_ServiceDesc is the grpc.ServiceDesc for VehicleData service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VehicleData_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "vehicledata.VehicleData",
	HandlerType: (*VehicleDataServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetVehicleDataByRegistration",
			Handler:    _VehicleData_GetVehicleDataByRegistration_Handler,
		},
		{
			MethodName: "GetVehicleDataByChassisEngine",
			Handler:    _VehicleData_GetVehicleDataByChassisEngine_Handler,
		},
		{
			MethodName: "GetDriverDataByLicenseNumber",
			Handler:    _VehicleData_GetDriverDataByLicenseNumber_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "vehicledata/vehicle-data.proto",
}
