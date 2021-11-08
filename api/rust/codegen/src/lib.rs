use proc_macro::TokenStream;
use quote::quote;
use std::iter::FromIterator;
use syn::{parse_macro_input, FnArg, ForeignItem};
use syn::{DeriveInput, ItemForeignMod};

#[proc_macro_derive(Runnable)]
pub fn derive_runnable(token_stream: TokenStream) -> TokenStream {
	let input = parse_macro_input!(token_stream as DeriveInput);

	let runnable_name = input.ident;

	let expanded = quote! {
		static RUNNABLE: &#runnable_name = &#runnable_name{};

		#[no_mangle]
		pub extern fn init() {
			suborbital::runnable::use_runnable(RUNNABLE);
		}
	};

	TokenStream::from(expanded)
}

fn create_function_wrapper(func: &ForeignItem) -> proc_macro2::TokenStream {
	match func {
		ForeignItem::Fn(func) => {
			let name = &func.sig.ident;
			let mut params = func.sig.inputs.clone();
			let attrs = proc_macro2::TokenStream::from_iter(func.attrs.iter().map(|attr| quote! {#attr}));
			params.pop();
			// // TODO: ensure error is returned
			let ident = quote! {crate::STATE.ident};
			let mut args_vec: Vec<proc_macro2::TokenStream> = params
				.iter()
				.map(|p| match p {
					FnArg::Typed(type_) => type_.pat.clone(),
					_ => panic!("Unexpected Type in ABI"),
				})
				.map(|p| quote! { #p,})
				.collect();
			args_vec.push(ident);
			let args = proc_macro2::TokenStream::from_iter(args_vec);
			let return_val = &func.sig.output;
			// Remove last arg
			quote! {
			#attrs
					pub fn #name (#params) #return_val {
				unsafe { super::#name(#args) }
					}
				  }
		}
		_ => quote! {},
	}
}

#[proc_macro]
pub fn wrap_host_functions(token_stream: TokenStream) -> TokenStream {
	if let Ok(extern_block) = syn::parse::<ItemForeignMod>(token_stream.clone()) {
		let funcs = proc_macro2::TokenStream::from_iter(extern_block.items.iter().map(create_function_wrapper));

		let expanded = quote! {
		#extern_block
		pub mod env {
		  #funcs
		}
		};
		TokenStream::from(expanded)
	} else {
		token_stream
	}
}
